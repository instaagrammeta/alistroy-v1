export type Locale = 'tg' | 'ru'
export type UserRole = 'customer' | 'seller' | 'driver' | 'admin'
export type ProductStatus = 'draft' | 'pending' | 'approved' | 'rejected'
export type ContactOwner = 'admin' | 'seller'
export type OrderStatus = 'new' | 'processing' | 'assigned' | 'on_delivery' | 'completed' | 'cancelled'
export type ReviewStatus = 'pending' | 'approved' | 'rejected'

export interface User {
  id: string
  email: string
  phone: string
  login: string
  name: string
  role: UserRole
  status: string
  locale: Locale
  avatar_url: string
  created_at: string
  customer?: Customer | null
  seller?: Seller | null
  driver?: Driver | null
}

export interface Customer {
  id: string
  user_id: string
  company: string
  address: string
  city: string
  notes: string
}

export interface Seller {
  id: string
  user_id: string
  full_name: string
  company_name: string
  market_name: string
  slug: string
  phone: string
  phone_alt: string
  whatsapp: string
  telegram: string
  telegram_username: string
  address: string
  city: string
  logo_url: string
  notes: string
  active: boolean
  is_featured: boolean
  created_at: string
}

export interface Driver {
  id: string
  user_id: string
  full_name: string
  age: number
  phone: string
  phone_alt: string
  whatsapp: string
  telegram: string
  vehicle: string
  photo_url: string
  notes: string
  active: boolean
  on_duty: boolean
}

export interface Category {
  id: string
  slug: string
  name_tj: string
  name_ru: string
  description_tj: string
  description_ru: string
  icon_url: string
  banner_url: string
  sort_order: number
  active: boolean
  subcategories?: Subcategory[]
}

export interface Subcategory {
  id: string
  category_id: string
  slug: string
  name_tj: string
  name_ru: string
  icon_url: string
  sort_order: number
  active: boolean
}

export interface Brand {
  id: string
  slug: string
  name: string
  logo_url: string
  active: boolean
  sort_order: number
}

export interface ProductImage {
  id: string
  product_id: string
  url: string
  alt: string
  sort_order: number
  is_cover: boolean
}

export interface Product {
  id: string
  seller_id: string
  category_id: string
  subcategory_id?: string | null
  brand_id?: string | null
  slug: string
  sku: string
  name_tj: string
  name_ru: string
  description_tj: string
  description_ru: string
  unit: string
  currency: string
  cost_price: number
  sale_price: number
  discount_percent: number
  stock_quantity: number
  minimum_stock: number
  is_available: boolean
  is_featured: boolean
  contact_owner: ContactOwner
  contact_phone: string
  contact_whatsapp: string
  contact_telegram: string
  status: ProductStatus
  rejection_note: string
  views_count: number
  phone_clicks: number
  whatsapp_clicks: number
  telegram_clicks: number
  created_at: string
  seller?: Seller | null
  category?: Category | null
  subcategory?: Subcategory | null
  brand?: Brand | null
  images?: ProductImage[]
}

export interface Banner {
  id: string
  position: string
  title_tj: string
  title_ru: string
  description_tj: string
  description_ru: string
  desktop_url: string
  tablet_url: string
  mobile_url: string
  link_url: string
  sort_order: number
  active: boolean
}

export interface Review {
  id: string
  product_id: string
  user_id: string
  rating: number
  comment: string
  status: ReviewStatus
  created_at: string
  user?: User | null
  product?: Product | null
}

export interface CartItem {
  id: string
  customer_id: string
  product_id: string
  quantity: number
  product?: Product | null
}

export interface Favorite {
  id: string
  user_id: string
  product_id: string
  product?: Product | null
}

export interface OrderItem {
  id: string
  order_id: string
  product_id: string
  name_snapshot: string
  unit: string
  quantity: number
  cost_price: number
  sale_price: number
  line_total: number
  profit: number
  product?: Product | null
}

export interface Order {
  id: string
  number: string
  customer_id?: string | null
  customer_name: string
  customer_phone: string
  delivery_address: string
  delivery_date?: string | null
  status: OrderStatus
  discount_percent: number
  subtotal: number
  total: number
  cost_total: number
  profit: number
  currency: string
  notes: string
  driver_id?: string | null
  created_at: string
  items?: OrderItem[]
  customer?: Customer | null
  driver?: Driver | null
}

export interface Notification {
  id: string
  user_id: string
  kind: string
  title_tj: string
  title_ru: string
  body_tj: string
  body_ru: string
  link_url: string
  read_at?: string | null
  created_at: string
}

export interface ChatRoom {
  id: string
  customer_id: string
  last_message_at?: string | null
  unread_admin: number
  unread_customer: number
  customer?: Customer | null
}

export interface ChatAttachment {
  id: string
  message_id: string
  url: string
  mime_type: string
  size_bytes: number
}

export interface ChatMessage {
  id: string
  room_id: string
  sender_id: string
  sender_role: string
  body: string
  read_at?: string | null
  created_at: string
  attachments?: ChatAttachment[]
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  access_expires_at: string
  refresh_expires_at: string
  user: User
}

export interface Pagination {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: Pagination
}

export interface AdminTotals {
  total_products: number
  approved_products: number
  pending_products: number
  total_categories: number
  total_brands: number
  total_sellers: number
  total_customers: number
  total_drivers: number
  total_orders: number
  new_orders: number
  total_views: number
  total_phone_clicks: number
  total_whatsapp_clicks: number
  total_telegram_clicks: number
  total_revenue: number
  total_profit: number
}

export interface SellerTotals {
  total_products: number
  approved_products: number
  pending_products: number
  low_stock: number
  total_views: number
  total_phone_clicks: number
  total_whatsapp_clicks: number
  total_telegram_clicks: number
}

export interface FinancialSummary {
  income: number
  expense: number
  purchase: number
  profit: number
}

export type SettingsMap = Record<string, string>

export interface BoardProduct { id: string; name_tj: string; name_ru: string; slug: string; price: number; status: string }
export interface BoardSubcategory { id: string; name_tj: string; name_ru: string; products: BoardProduct[] }
export interface BoardCategory { id: string; name_tj: string; name_ru: string; slug: string; subcategories: BoardSubcategory[]; products: BoardProduct[] }
