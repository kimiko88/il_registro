-- Migration 132: Elezioni Organi Collegiali (Urna Digitale & Voto Online Anonimo)
-- Normativa: O.M. 215/1991, Linee guida voto elettronico, GDPR Art. 5 (anonimato del voto)

CREATE TABLE IF NOT EXISTS school_elections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    election_tier VARCHAR(30) NOT NULL CHECK (election_tier IN ('classe', 'consiglio_istituto', 'consulta_studenti')),
    target_role VARCHAR(20) NOT NULL CHECK (target_role IN ('parent', 'student', 'teacher', 'ata')),
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    max_preferences INT NOT NULL DEFAULT 1,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_school_elections_school_tier
    ON school_elections(school_id, election_tier);

CREATE TABLE IF NOT EXISTS election_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    election_id UUID NOT NULL REFERENCES school_elections(id) ON DELETE CASCADE,
    list_number INT NOT NULL,
    motto VARCHAR(255) NOT NULL,
    UNIQUE(election_id, list_number)
);

CREATE TABLE IF NOT EXISTS election_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL REFERENCES election_lists(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    candidate_order INT NOT NULL,
    UNIQUE(list_id, candidate_order)
);

-- Registro elettori: traccia chi ha votato (senza salvare la preferenza espressa!)
CREATE TABLE IF NOT EXISTS election_voter_registry (
    election_id UUID NOT NULL REFERENCES school_elections(id) ON DELETE CASCADE,
    voter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    voted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    receipt_token VARCHAR(64) NOT NULL,
    PRIMARY KEY(election_id, voter_id)
);

-- Urna anonima disaccoppiata (nessun link con voter_id!)
CREATE TABLE IF NOT EXISTS election_ballot_box (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    election_id UUID NOT NULL REFERENCES school_elections(id) ON DELETE CASCADE,
    list_id UUID REFERENCES election_lists(id) ON DELETE CASCADE,
    candidate_ids UUID[] DEFAULT ARRAY[]::UUID[],
    is_blank BOOLEAN NOT NULL DEFAULT FALSE,
    cast_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_election_ballot_box_election
    ON election_ballot_box(election_id);

-- RLS
ALTER TABLE school_elections ENABLE ROW LEVEL SECURITY;
ALTER TABLE election_lists ENABLE ROW LEVEL SECURITY;
ALTER TABLE election_candidates ENABLE ROW LEVEL SECURITY;
ALTER TABLE election_voter_registry ENABLE ROW LEVEL SECURITY;
ALTER TABLE election_ballot_box ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    DROP POLICY IF EXISTS "School members can view elections" ON school_elections;
    CREATE POLICY "School members can view elections"
        ON school_elections FOR SELECT
        USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

    DROP POLICY IF EXISTS "Admin and secretary can manage elections" ON school_elections;
    CREATE POLICY "Admin and secretary can manage elections"
        ON school_elections FOR ALL
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = school_elections.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
        ));

    DROP POLICY IF EXISTS "School members can view lists and candidates" ON election_lists;
    DROP POLICY IF EXISTS election_lists_select_policy ON election_lists;
    DROP POLICY IF EXISTS election_lists_service_policy ON election_lists;
    CREATE POLICY election_lists_select_policy
        ON election_lists FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_lists_service_policy
        ON election_lists FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    DROP POLICY IF EXISTS "School members can view candidates" ON election_candidates;
    DROP POLICY IF EXISTS election_candidates_select_policy ON election_candidates;
    DROP POLICY IF EXISTS election_candidates_service_policy ON election_candidates;
    CREATE POLICY election_candidates_select_policy
        ON election_candidates FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_candidates_service_policy
        ON election_candidates FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    DROP POLICY IF EXISTS "Voters can check own registry" ON election_voter_registry;
    DROP POLICY IF EXISTS election_voter_registry_select_policy ON election_voter_registry;
    DROP POLICY IF EXISTS election_voter_registry_service_policy ON election_voter_registry;
    CREATE POLICY election_voter_registry_select_policy
        ON election_voter_registry FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_voter_registry_service_policy
        ON election_voter_registry FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    DROP POLICY IF EXISTS "Ballot box anonymous insert and count" ON election_ballot_box;
    DROP POLICY IF EXISTS election_ballot_box_select_policy ON election_ballot_box;
    DROP POLICY IF EXISTS election_ballot_box_service_policy ON election_ballot_box;
    CREATE POLICY election_ballot_box_select_policy
        ON election_ballot_box FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_ballot_box_service_policy
        ON election_ballot_box FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);
END $$;
