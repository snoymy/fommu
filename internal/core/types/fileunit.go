package types

type Byte uint64

func KiB(kib uint64) Byte {
    return Byte(kib*1024)
}

func MiB(mib uint64) Byte {
    return KiB(mib*1024)
}

func GiB(gib uint64) Byte {
    return MiB(gib*1024)
}

