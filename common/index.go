package common

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"
	"path"
	"path/filepath"
)

const (
	INDEX = "index"
)

// IndexEntry 文件中均二进制存储
type IndexEntry struct {
	// 0-7 8字节
	// ctime_sec 4字节 ctime_nsec 4字节
	Ctime int64
	// 8-15 8字节
	// mtime_sec 4字节 mtime_nsec 4字节
	Mtime int64
	// 16-19 4字节
	Dev uint32
	// 20-23 4字节
	Ino uint32
	// 24-25 2字节 未使用
	// 26-27 2字节 mode
	// ModeType 为 mode >> 12
	ModeType uint16
	// ModePerms 为 mode & 0b0000000111111111
	ModePerms uint16
	// 28-31 4字节
	Uid uint32
	// 32-35 4字节
	Gid uint32
	// 36-39 4字节
	Size uint32
	// 40-59 20字节
	Sha1 string
	// 60-61 2字节 flags
	// FlagAssumeValid 为 flags & 0b1000000000000000 != 0
	// FlagExtended 为 (flags & 0b0100000000000000) != 0 但是此处未使用
	FlagAssumeValid bool
	// FlagStage 为 flags & 0b0011000000000000
	FlagStage uint16
	// Name 长度不定 nameLength = flags & 0b0000111111111111
	Name string
}

type GitIndex struct {
	Version int
	Entries []*IndexEntry
}

func (index *GitIndex) AddEntry(entry *IndexEntry) {
	index.Entries = append(index.Entries, entry)
}

func (index *GitIndex) UpdateEntry(entry *IndexEntry) {
	for j := range index.Entries {
		if index.Entries[j].Name == entry.Name {
			index.Entries[j] = entry
			return
		}
	}
}

func ReadIndex() (*GitIndex, error) {
	data, err := os.ReadFile(path.Join(CURRENT, REPO, INDEX))
	if err != nil {
		return nil, errors.New("error reading index file: " + err.Error())
	}
	// header
	header := data[:12]
	signature := header[:4]
	if string(signature) != "DIRC" {
		return nil, errors.New("invalid index file signature: " + string(signature))
	}
	version := int(binary.BigEndian.Uint32(data[4:8]))
	count := int(binary.BigEndian.Uint32(data[8:12]))
	// entries
	entries := make([]*IndexEntry, count)
	content := data[12:]
	idx := 0
	for i := 0; i < count; i++ {
		entries[i] = &IndexEntry{
			Ctime:           int64(binary.BigEndian.Uint64(content[idx : idx+8])),
			Mtime:           int64(binary.BigEndian.Uint64(content[idx+8 : idx+16])),
			Dev:             uint32(binary.BigEndian.Uint32(content[idx+16 : idx+20])),
			Ino:             uint32(binary.BigEndian.Uint32(content[idx+20 : idx+24])),
			ModeType:        uint16(binary.BigEndian.Uint16(content[idx+26:idx+28]) >> 12),
			ModePerms:       uint16(binary.BigEndian.Uint16(content[idx+26:idx+28]) & 0b0000000111111111),
			Uid:             uint32(binary.BigEndian.Uint32(content[idx+28 : idx+32])),
			Gid:             uint32(binary.BigEndian.Uint32(content[idx+32 : idx+36])),
			Size:            uint32(binary.BigEndian.Uint32(content[idx+36 : idx+40])),
			Sha1:            hex.EncodeToString(content[idx+40 : idx+60]),
			FlagAssumeValid: binary.BigEndian.Uint16(content[idx+60:idx+62])&0b1000000000000000 != 0,
			FlagStage:       uint16(binary.BigEndian.Uint16(content[idx+60:idx+62]) & 0b0011000000000000),
		}
		nameLength := binary.BigEndian.Uint16(content[idx+60:idx+62]) & 0b0000111111111111
		idx += 62
		rawName := content[idx : idx+int(nameLength)]
		idx += int(nameLength) + 1

		entries[i].Name = string(rawName)

		// idx = 8 * ceil(idx / 8)
		idx = (idx + 7) &^ 7
	}
	return &GitIndex{
		Version: version,
		Entries: entries,
	}, nil
}

func (index *GitIndex) WriteIndex() error {
	data := make([]byte, 0)
	// header
	header := make([]byte, 12)
	copy(header[:4], "DIRC")
	binary.BigEndian.PutUint32(header[4:8], uint32(index.Version))
	binary.BigEndian.PutUint32(header[8:12], uint32(len(index.Entries)))
	data = append(data, header...)
	// entries
	idx := 0
	for i := range index.Entries {
		entry, err := encodeIndexEntry(index.Entries[i])
		if err != nil {
			return err
		}
		data = append(data, entry...)
		idx += 62 + len(index.Entries[i].Name) + 1
		//  add padding if necessary.
		if idx%8 != 0 {
			padSize := 8 - (idx % 8)
			pad := make([]byte, padSize)
			data = append(data, pad...)
			idx += padSize
		}
	}
	return os.WriteFile(path.Join(CURRENT, REPO, INDEX), data, 0644)
}

func encodeIndexEntry(entry *IndexEntry) ([]byte, error) {
	data := make([]byte, 62)
	binary.BigEndian.PutUint64(data[0:8], uint64(entry.Ctime))
	binary.BigEndian.PutUint64(data[8:16], uint64(entry.Mtime))
	binary.BigEndian.PutUint32(data[16:20], entry.Dev)
	binary.BigEndian.PutUint32(data[20:24], entry.Ino)
	mode := (entry.ModeType << 12) | entry.ModePerms
	binary.BigEndian.PutUint16(data[26:28], mode)
	binary.BigEndian.PutUint32(data[28:32], entry.Uid)
	binary.BigEndian.PutUint32(data[32:36], entry.Gid)
	binary.BigEndian.PutUint32(data[36:40], entry.Size)
	sha1, err := hex.DecodeString(entry.Sha1)
	if err != nil {
		return nil, err
	}
	copy(data[40:60], sha1)
	flags := uint16(0)
	if entry.FlagAssumeValid {
		flags |= 0b1000000000000000
	}
	flags |= entry.FlagStage
	nameLength := uint16(len(entry.Name))
	flags |= nameLength
	binary.BigEndian.PutUint16(data[60:62], flags)
	data = append(data, []byte(entry.Name)...)
	data = append(data, 0x00)
	return data, nil
}

func NewIndexEntry(name string) (*IndexEntry, error) {
	stat, err := os.Stat(filepath.Join(CURRENT, name))
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}
	sha1, err := GetSha1ByPath("blob", name, true)
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}
	return &IndexEntry{
		Ctime:           stat.ModTime().Unix(),
		Mtime:           stat.ModTime().Unix(),
		Dev:             0,
		Ino:             0,
		ModeType:        0b1000,
		ModePerms:       0o644,
		Uid:             0,
		Gid:             0,
		Size:            uint32(stat.Size()),
		Sha1:            sha1,
		FlagAssumeValid: false,
		FlagStage:       0,
		Name:            name,
	}, nil
}

func GetIndexInfo() (map[string]string, error) {
	index, err := ReadIndex()
	if err != nil {
		return nil, errors.New("error reading index: " + err.Error())
	}
	indexNameToSha1 := make(map[string]string)

	for i := range index.Entries {
		indexNameToSha1[index.Entries[i].Name] = index.Entries[i].Sha1
	}
	return indexNameToSha1, nil
}

func (index *GitIndex) MakeTreeContent() ([]byte, error) {
	gitIndex, err := ReadIndex()
	if err != nil {
		return nil, errors.New("error reading index: " + err.Error())
	}
	content := make([]byte, 0)
	for i := range gitIndex.Entries {
		content = append(content, "100644 "...)
		content = append(content, gitIndex.Entries[i].Name...)
		content = append(content, 0x00)
		sha1, _ := hex.DecodeString(gitIndex.Entries[i].Sha1)
		content = append(content, sha1...)
	}
	return content, nil
}
