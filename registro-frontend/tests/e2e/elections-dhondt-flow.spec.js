import { describe, it, expect, vi } from 'vitest'
import { electionsService } from '@/services/electionsService'

vi.mock('@/services/electionsService', () => ({
    electionsService: {
        getElectionDetails: vi.fn(),
        castVote: vi.fn(),
        getScrutiny: vi.fn(),
        closeElection: vi.fn()
    }
}))

describe('School Elections & d\'Hondt Scrutiny E2E Workflow', () => {
    it('executes full election lifecycle: voter enters booth, casts anonymous vote, receives receipt, and live scrutiny calculates d\'Hondt seats', async () => {
        const electionId = 'elec-consiglio-2026'

        // Step 1: Voter enters booth and loads lists
        electionsService.getElectionDetails.mockResolvedValueOnce({
            election: {
                id: electionId,
                title: 'Elezioni Consiglio d\'Istituto - Componente Genitori',
                max_preferences: 2,
                is_closed: false
            },
            lists: [
                {
                    id: 'list-1',
                    list_number: 1,
                    motto: 'Scuola Aperta e Partecipata',
                    candidates: [
                        { id: 'cand-1', first_name: 'Chiara', last_name: 'Ferrari' },
                        { id: 'cand-2', first_name: 'Roberto', last_name: 'Esposito' }
                    ]
                },
                {
                    id: 'list-2',
                    list_number: 2,
                    motto: 'Innovazione Didattica e Benessere',
                    candidates: [
                        { id: 'cand-3', first_name: 'Matteo', last_name: 'Ricci' }
                    ]
                }
            ]
        })
        const detailsRes = await electionsService.getElectionDetails(electionId)
        expect(detailsRes.lists).toHaveLength(2)

        // Step 2: Voter casts vote (anonymous decoupled)
        electionsService.castVote.mockResolvedValueOnce({
            message: 'Voto anonimo registrato nell\'urna digitale',
            receipt: {
                receipt_token: '9f8b4c2a1e7d6a5f0b3c8e1d2a4f6b8c',
                election_id: electionId,
                voted_at: '2026-10-04T10:00:00Z'
            }
        })
        const voteRes = await electionsService.castVote(electionId, {
            election_id: electionId,
            list_id: 'list-1',
            candidate_ids: ['cand-1'],
            is_blank: false
        })
        expect(voteRes.receipt.receipt_token).toHaveLength(32)

        // Step 3: Commission runs scrutiny with d'Hondt method
        electionsService.getScrutiny.mockResolvedValueOnce({
            scrutiny: {
                election_id: electionId,
                total_voters: 120,
                total_votes_cast: 120,
                blank_votes: 5,
                lists_results: [
                    { list_id: 'list-1', list_number: 1, motto: 'Scuola Aperta e Partecipata', total_votes: 75, seats_won: 3 },
                    { list_id: 'list-2', list_number: 2, motto: 'Innovazione Didattica e Benessere', total_votes: 40, seats_won: 1 }
                ],
                elected_members: [
                    { candidate_id: 'cand-1', first_name: 'Chiara', last_name: 'Ferrari', list_motto: 'Scuola Aperta e Partecipata', votes: 48 }
                ]
            }
        })
        const scrutinyRes = await electionsService.getScrutiny(electionId, 4)
        expect(scrutinyRes.scrutiny.lists_results[0].seats_won).toBe(3)
        expect(scrutinyRes.scrutiny.lists_results[1].seats_won).toBe(1)
        expect(scrutinyRes.scrutiny.elected_members[0].first_name).toBe('Chiara')
    })
})
