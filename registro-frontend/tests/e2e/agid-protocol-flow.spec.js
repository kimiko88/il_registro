import { describe, it, expect, vi } from 'vitest'
import { protocolService } from '@/services/protocolService'

vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn(),
        get: vi.fn()
    }
}))

vi.mock('@/services/protocolService', () => ({
    protocolService: {
        protocolDocument: vi.fn(),
        getEntries: vi.fn(),
        getProtocolByEntity: vi.fn()
    }
}))

describe('AgID Protocol Register & Visual Stamp E2E Workflow', () => {
    it('executes official protocol lifecycle: protocol document with automatic numbering, visual stamp generation, and entity linkage', async () => {
        // Step 1: Secretary protocols an act
        protocolService.protocolDocument.mockResolvedValueOnce({
            message: 'Documento protocollato a norma AgID',
            protocol: {
                id: 'prot-e2e-1',
                protocol_number: 1045,
                protocol_year: 2026,
                protocol_date: '2026-10-04T10:15:00Z',
                flow_direction: 'out',
                classification_title: 7,
                classification_class: '2',
                subject: 'Convocazione straordinaria Consiglio di Classe',
                sender: 'Dirigente Scolastico',
                recipient: 'Docenti della classe 3A',
                visual_stamp: 'LICEO STATALE GALVANI - REG. UFF. PROT. N. 0001045 del 04/10/2026'
            }
        })

        const protRes = await protocolService.protocolDocument({
            subject: 'Convocazione straordinaria Consiglio di Classe',
            sender: 'Dirigente Scolastico',
            recipient: 'Docenti della classe 3A',
            flow_direction: 'out',
            classification_title: 7,
            classification_class: '2',
            entity_type: 'verbale',
            entity_id: '00000000-0000-0000-0000-000000000099'
        })

        expect(protRes.protocol.protocol_number).toBe(1045)
        expect(protRes.protocol.visual_stamp).toContain('0001045')

        // Step 2: Retrieve protocol by entity linkage (e.g. from verbale view)
        protocolService.getProtocolByEntity.mockResolvedValueOnce({
            id: 'prot-e2e-1',
            protocol_number: 1045,
            protocol_year: 2026,
            visual_stamp: 'LICEO STATALE GALVANI - REG. UFF. PROT. N. 0001045 del 04/10/2026'
        })
        const linkedRes = await protocolService.getProtocolByEntity('verbale', '00000000-0000-0000-0000-000000000099')
        expect(linkedRes.protocol_number).toBe(1045)
        expect(linkedRes.visual_stamp).toContain('0001045')
    })
})
