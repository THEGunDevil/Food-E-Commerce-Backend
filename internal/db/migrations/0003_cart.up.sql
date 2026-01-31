-- +migrate Up

CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID UNIQUE NOT NULL,

    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'checked_out', 'expired')),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    cart_id UUID NOT NULL
        REFERENCES carts(id)
        ON DELETE CASCADE,

    menu_item_id UUID NOT NULL
        REFERENCES menu_items(id)
        ON DELETE CASCADE,

    quantity INTEGER NOT NULL
        DEFAULT 1
        CHECK (quantity >= 1 AND quantity <= 10),

    special_instructions TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (cart_id, menu_item_id)
);

