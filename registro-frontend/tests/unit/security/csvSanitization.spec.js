import { describe, it, expect } from 'vitest'

describe('Security: CSV Formula Injection Defense', () => {
    const sanitizeCSVCell = (val) => {
        if (val === null || val === undefined) return '""'
        let str = String(val)
        if (/^[=+\-@\t\r]/.test(str)) {
            str = "'" + str
        }
        return `"${str.replace(/"/g, '""')}"`
    }

    it('should prefix leading formula characters (=, +, -, @) with single quote', () => {
        expect(sanitizeCSVCell("=CMD|'/C calc'!A1")).toBe(`"'=CMD|'/C calc'!A1"`)
        expect(sanitizeCSVCell('+SUM(1+1)')).toBe(`"'+SUM(1+1)"`)
        expect(sanitizeCSVCell('-100')).toBe(`"'-100"`)
        expect(sanitizeCSVCell('@malicious')).toBe(`"'@malicious"`)
    })

    it('should properly escape double quotes and handle safe text', () => {
        expect(sanitizeCSVCell('Mario Rossi')).toBe(`"Mario Rossi"`)
        expect(sanitizeCSVCell('John "The Boss" Doe')).toBe(`"John ""The Boss"" Doe"`)
    })
})
