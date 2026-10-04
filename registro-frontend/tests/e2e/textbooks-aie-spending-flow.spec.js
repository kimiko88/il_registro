import { describe, it, expect, vi } from 'vitest'
import { textbookService } from '@/services/textbookService'

vi.mock('@/services/textbookService', () => ({
    textbookService: {
        importAIE: vi.fn(),
        getSpendingReport: vi.fn(),
        adoptBook: vi.fn(),
        exportClassAIE: vi.fn()
    }
}))

describe('AIE Textbooks & Spending Limits E2E Workflow', () => {
    it('executes full cycle: import AIE catalog, compute spending status, check tolerance threshold, and export ministerial txt', async () => {
        // Step 1: Simulate AIE Catalog Import
        const formData = new FormData()
        textbookService.importAIE.mockResolvedValueOnce({
            data: { message: 'Catalogo AIE importato con successo', imported: 2 }
        })
        const importRes = await textbookService.importAIE(formData)
        expect(importRes.data.imported).toBe(2)

        // Step 2: Simulate Class Adoption
        textbookService.adoptBook.mockResolvedValueOnce({ status: 201 })
        const adoptRes = await textbookService.adoptBook('class-1', {
            book_id: 'book-1',
            subject_id: 'subj-math',
            adoption_type: 'nuova_adozione'
        })
        expect(adoptRes.status).toBe(201)

        // Step 3: Check Spending Report with Tolerance (315€ over 300€ limit but under 330€ max tolerance)
        textbookService.getSpendingReport.mockResolvedValueOnce({
            data: {
                total_spending: 315.00,
                spending_limit: 300.00,
                tolerance_threshold: 330.00,
                difference: 15.00,
                percentage: 105.00,
                status: 'WARNING_TOLERANCE'
            }
        })
        const reportRes = await textbookService.getSpendingReport('class-1')
        expect(reportRes.data.status).toBe('WARNING_TOLERANCE')
        expect(reportRes.data.difference).toBe(15.00)

        // Step 4: Export Ministerial AIE .txt
        textbookService.exportClassAIE.mockResolvedValueOnce({
            data: 'RMPS010004;2026/2027;1A;MAT;9788808836243;Matematica.blu 2.0;Bergamini;Zanichelli;28.90;1;nuova_adozione;NO'
        })
        const exportRes = await textbookService.exportClassAIE('class-1')
        expect(exportRes.data).toContain('RMPS010004')
        expect(exportRes.data).toContain('9788808836243')
    })
})
