-- name: CreateMenuItem :one
INSERT INTO menu_items (
    category_id,
    name,
    slug,
    description,
    price,
    discount_price,
    ingredients,
    tags,
    prep_time,
    spicy_level,
    is_vegetarian,
    is_special,
    is_available,
    stock_quantity,
    min_stock_alert,
    display_order
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16
)
RETURNING *;

-- name: CreateMenuItemImage :exec
INSERT INTO menu_item_images (
    menu_item_id,
    image_url,
    image_public_id,
    is_primary,
    display_order
)
VALUES ($1, $2, $3, $4, $5);

-- name: GetMenuItemImagesByMenuItemID :many
SELECT *
FROM menu_item_images
WHERE menu_item_id = $1;

-- name: ListMenuItemImages :many
SELECT *
FROM menu_item_images
WHERE menu_item_id = $1
ORDER BY is_primary DESC, display_order ASC;

-- name: SetPrimaryMenuItemImage :exec
UPDATE menu_item_images
SET is_primary = CASE
    WHEN id = $2 THEN true
    ELSE false
END
WHERE menu_item_id = $1;


-- name: DeleteMenuItemImage :exec
DELETE FROM menu_item_images
WHERE id = $1;

-- name: GetMenuItem :one
SELECT * FROM menu_items 
WHERE id = $1 LIMIT 1;

-- name: GetMenuItemByID :one
SELECT
    mi.*,
    c.name AS CategoryName
FROM menu_items mi
JOIN categories c
    ON mi.category_id = c.id
WHERE mi.id = $1;



-- name: GetMenuItemBySlug :one
SELECT * FROM menu_items 
WHERE slug = $1 LIMIT 1;

-- name: ListMenuItems :many
SELECT
    menu_items.*,
    categories.name AS CategoryName
FROM menu_items
JOIN categories ON menu_items.category_id = categories.id
ORDER BY
    menu_items.display_order,
    menu_items.name
LIMIT $1 OFFSET $2;

-- name: ListMenuItemsByCategory :many
SELECT
    menu_items.*,
    categories.name AS CategoryName
FROM menu_items
JOIN categories ON menu_items.category_id = categories.id
WHERE category_id = $1 
AND is_available = true 
ORDER BY menu_items.display_order, menu_items.name 
LIMIT $2 OFFSET $3;

-- name: ListAvailableMenuItems :many
SELECT * FROM menu_items 
WHERE is_available = true 
ORDER BY display_order, name 
LIMIT $1 OFFSET $2;



-- name: ListVegetarianMenuItems :many
SELECT * FROM menu_items 
WHERE is_vegetarian = true 
AND is_available = true 
ORDER BY display_order, name 
LIMIT $1 OFFSET $2;

-- name: ListSpecialMenuItems :many
SELECT * FROM menu_items 
WHERE is_special = true 
AND is_available = true 
ORDER BY display_order, name 
LIMIT $1 OFFSET $2;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET
    category_id = COALESCE($2, category_id),
    name = COALESCE($3, name),
    slug = COALESCE($4, slug),
    description = COALESCE($5, description),
    price = COALESCE($6, price),
    discount_price = COALESCE($7, discount_price),
    ingredients = COALESCE($8, ingredients),
    tags = COALESCE($9, tags),
    prep_time = COALESCE($10, prep_time),
    spicy_level = COALESCE($11, spicy_level),
    is_vegetarian = COALESCE($12, is_vegetarian),
    is_special = COALESCE($13, is_special),
    is_available = COALESCE($14, is_available),
    stock_quantity = COALESCE($15, stock_quantity),
    min_stock_alert = COALESCE($16, min_stock_alert),
    display_order = COALESCE($17, display_order),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: DeleteMenuItem :exec
DELETE FROM menu_items 
WHERE id = $1;

-- name: ToggleMenuItemAvailability :one
UPDATE menu_items 
SET 
    is_available = NOT is_available,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateMenuItemStock :one
UPDATE menu_items 
SET 
    stock_quantity = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: IncrementMenuItemOrders :one
UPDATE menu_items 
SET 
    total_orders = total_orders + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateMenuItemRating :one
UPDATE menu_items 
SET 
    average_rating = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: SearchMenuItems :many
SELECT * FROM menu_items 
WHERE 
    (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND is_available = true
ORDER BY display_order, name 
LIMIT $2 OFFSET $3;

-- name: FilterMenuItemsByTags :many
SELECT * FROM menu_items 
WHERE 
    tags @> $1::varchar(50)[]
    AND is_available = true
ORDER BY display_order, name 
LIMIT $2 OFFSET $3;

-- name: FilterMenuItemsBySpicyLevel :many
SELECT * FROM menu_items 
WHERE 
    spicy_level = $1
    AND is_available = true
ORDER BY display_order, name 
LIMIT $2 OFFSET $3;

-- name: GetLowStockMenuItems :many
SELECT * FROM menu_items 
WHERE 
    stock_quantity > 0 
    AND stock_quantity <= min_stock_alert
    AND is_available = true
ORDER BY stock_quantity ASC;

-- name: GetOutOfStockMenuItems :many
SELECT * FROM menu_items 
WHERE 
    stock_quantity = 0
    AND is_available = true
ORDER BY name;

-- name: UpdateDisplayOrder :exec
UPDATE menu_items 
SET 
    display_order = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: BulkUpdateAvailability :exec
UPDATE menu_items 
SET 
    is_available = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE category_id = $1;

-- name: CountMenuItems :one
SELECT COUNT(*) FROM menu_items;
-- name: CountMenuItemsByCategory :one
SELECT COUNT(*) 
FROM menu_items
WHERE category_id = $1;

-- name: CountAvailableMenuItems :one
SELECT COUNT(*) FROM menu_items 
WHERE is_available = true;

-- name: IncrementStock :one
UPDATE menu_items 
SET 
    stock_quantity = stock_quantity + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DecrementStock :one
UPDATE menu_items 
SET 
    stock_quantity = stock_quantity - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetMenuItemsWithCategory :many
SELECT 
    m.*,
    c.name as category_name,
    c.slug as category_slug
FROM menu_items m
LEFT JOIN categories c ON m.category_id = c.id
WHERE m.is_available = true
ORDER BY m.display_order, m.name
LIMIT $1 OFFSET $2;

-- name: GetMenuItemStats :one
SELECT 
    COUNT(*) as total_items,
    SUM(total_orders) as total_all_orders,
    AVG(price) as avg_price,
    MIN(price) as min_price,
    MAX(price) as max_price
FROM menu_items 
WHERE is_available = true;

-- name: GetPopularMenuItems :many
SELECT * FROM menu_items 
WHERE is_available = true
ORDER BY total_orders DESC, average_rating DESC
LIMIT $1;

-- name: GetDiscountedMenuItems :many
SELECT * FROM menu_items 
WHERE discount_price IS NOT NULL 
AND discount_price > 0
AND is_available = true
ORDER BY (price - discount_price) DESC
LIMIT $1 OFFSET $2;

-- name: UpdateMenuItemPartial :one
UPDATE menu_items 
SET 
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- This query uses jsonb to update dynamic fields
-- name: UpdateMenuItemDynamic :one
UPDATE menu_items 
SET 
    name = COALESCE(sqlc.narg('name'), name),
    price = COALESCE(sqlc.narg('price'), price),
    is_available = COALESCE(sqlc.narg('is_available'), is_available),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;