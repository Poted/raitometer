-- Tabela główna dla kalkulatora, powiązana 1-do-1 z projektem
CREATE TABLE calculators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID UNIQUE NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabela dla modułów w ramach kalkulatora
CREATE TABLE calculator_modules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calculator_id UUID NOT NULL REFERENCES calculators(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabela dla pojedynczych pozycji w ramach modułu (rozbudowana)
CREATE TABLE calculator_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_id UUID NOT NULL REFERENCES calculator_modules(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    
    -- Typ wyceny: 'per_unit' (za sztukę), 'per_sqm' (za m2)
    price_type VARCHAR(50) NOT NULL, 
    
    -- Klucz do dopasowania z odpowiedzią AI
    ai_tag VARCHAR(100) UNIQUE, 
    
    unit_price NUMERIC(10, 2) NOT NULL,
    
    -- Pola, które wypełni klient lub wykonawca, jeśli AI nie jest używane
    unit VARCHAR(50), 
    quantity NUMERIC(10, 2) NOT NULL DEFAULT 1.0,

    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);