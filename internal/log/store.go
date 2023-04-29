package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var enc = binary.BigEndian // エンディアン

const lenWidth = 8 // レコードの長さを保存するために使うバイト数

// ファイルを保持
// ファイルへの書き込みと読み込みを行う
type store struct {
	*os.File               // ファイルを埋め込む
	mu       sync.Mutex    // ロック
	buf      *bufio.Writer // バッファ
	size     uint64        // ファイルのサイズ
}

func newStore(f *os.File) (*store, error) {
	// ファイルのサイズを取得
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

/**
 * ファイルにレコードを追加する
 * 書き込まれたバイト数、書き込んだデータのオフセット、エラーを返す
 */
func (s *store) Append(data []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock() // ロック
	defer s.mu.Unlock()

	pos = s.size // ファイルのサイズを取得

	// ファイルに書き込むデータの長さを先頭に書き込む
	if err := binary.Write(s.buf, enc, uint64(len(data))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(data)
	if err != nil {
		return 0, 0, err
	}

	// ファイルのサイズを更新
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

// オフセットを指定してファイルからデータを読み込む
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock() // ロック
	defer s.mu.Unlock()

	// バッファをメモリに書き込む
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// レコードの長さを取得
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	// レコードの長さ分のスライスを作成
	// レコードを読み込む
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// オフセットを指定してファイルからデータを読み込む
// 第一引数で指定した配列の長さ分だけ読み込む
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// 永続化
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
