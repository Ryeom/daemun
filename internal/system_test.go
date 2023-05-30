package internal

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

type TT struct {
	B  uint8 // is a byte
	I  int   // it is int32 on my x86 32 bit PC
	P  *int  // it is int32 on my x86 32 bit PC
	S  string
	SS []string
}

func TestSystem(t *testing.T) {
	PrintMemUsage()

	var overall [][]int
	for i := 0; i < 4; i++ {
		a := make([]int, 0, 999999)
		overall = append(overall, a)
		PrintMemUsage()
		time.Sleep(time.Second)
	}

	overall = nil
	PrintMemUsage()

	runtime.GC()
	PrintMemUsage()

	const PtrSize = 32 << uintptr(^uintptr(0)>>63)
	fmt.Println("PtrSize=", PtrSize)
	fmt.Println("IntSize=", strconv.IntSize)
	var m1, m2, m3, m4, m5, m6 runtime.MemStats
	runtime.ReadMemStats(&m1)
	tt := TT{}
	runtime.ReadMemStats(&m2)
	fmt.Println("sizeof(uint8)", unsafe.Sizeof(tt.B), "offset=", unsafe.Offsetof(tt.B))
	fmt.Println("sizeof(int)", unsafe.Sizeof(tt.I), "offset=", unsafe.Offsetof(tt.I))
	fmt.Println("sizeof(*int)", unsafe.Sizeof(tt.P), "offset=", unsafe.Offsetof(tt.P))
	fmt.Println("sizeof(string)", unsafe.Sizeof(tt.S), "offset=", unsafe.Offsetof(tt.S))
	fmt.Println("sizeof([]string)", unsafe.Sizeof(tt.SS), "offset=", unsafe.Offsetof(tt.SS))
	fmt.Println("sizeof(T)", unsafe.Sizeof(tt))

	memUsage(&m1, &m2)

	runtime.ReadMemStats(&m3)
	t2 := "abc"
	runtime.ReadMemStats(&m4)
	memUsage(&m3, &m4)

	runtime.ReadMemStats(&m5)
	t3 := map[int]string{1: "x"}
	runtime.ReadMemStats(&m6)
	memUsage(&m5, &m6)
	fmt.Println(t2, t3)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func memUsage(m1, m2 *runtime.MemStats) {
	fmt.Println("Alloc:", m2.Alloc-m1.Alloc, "TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc, "HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc)
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 사용 중인 힙 메모리
	usedHeap := m.HeapAlloc / 1024 / 1024
	// 할당된 힙 메모리
	allocatedHeap := m.HeapSys / 1024 / 1024
	// 시스템 전체 메모리
	totalSys := m.Sys / 1024 / 1024

	fmt.Printf("Memory Usage: Used Heap %dMB / Allocated Heap %dMB / Total System %dMB\n", usedHeap, allocatedHeap, totalSys)
}
