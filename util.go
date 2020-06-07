package gograph

import "image/color"

type RGBA16 struct {
    R, G, B, A uint16
}

func (c RGBA16) RGBA() (r, g, b, a uint32) {
    r = uint32(float64(c.R) * float64(c.A) / 0xFFFF)
    g = uint32(float64(c.G) * float64(c.A) / 0xFFFF)
    b = uint32(float64(c.B) * float64(c.A) / 0xFFFF)
    a = uint32(c.A)

    return
}

func MinInt(a, b int) int {
    if a <= b {
        return a
    }

    return b
}

func BlendColor(old, new color.Color) color.Color {
    old_r, old_g, old_b, _ := old.RGBA()
    new_r, new_g, new_b, new_a := new.RGBA()

    return RGBA16{
        uint16((new_a * new_r + (0xFFFF - new_a) * old_r) / 0xFFFF),
        uint16((new_a * new_g + (0xFFFF - new_a) * old_g) / 0xFFFF),
        uint16((new_a * new_b + (0xFFFF - new_a) * old_b) / 0xFFFF),
        0xFFFF,
    }
}