// Shared API types between frontend and backend.

export type Locale = 'tg' | 'ru'
export type UserRole = 'customer' | 'seller' | 'admin'
export type ProductStatus = 'draft' | 'pending' | 'approved' | 'rejected'
export type ContactType = 'admin' | 'seller'
export type SellerStatus = 'pending' | 'approved' | 'blocked'
export type ReviewStatus = 'pending' | 'approved' | 'rejected'

export interface User {
  id: string
  email: string
  name: string
  phone: string
  role: UserRole
  locale: Locale
  is_active: boolean
  created_at: string
  updated_at: string
  seller?: Seller | null
}

export interface Seller {
  id: string
  user_id: string
  name: string
  slug: string
  description_tj: string
  description_ru: string
  logo_url: string
  phone: string
  whatsapp: string
  address: string
  city: string
  status: SellerStatus
  is_featured: boolean
  created_at: string
  updated_at: string
  user?: User | null
}

export interface Category {
  id: string
  slug: string
  title_tj: string
  title_ru: string
  icon_url: string
  sort_order: number
  is_active: boolean
  parent_id?: string | null
  created_at: string
  updated_at: string
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
  slug: string
  sku: string
  title_tj: string
  title_ru: string
  description_tj: string
  description_ru: string
  price: number
  currency: string
  unit: string
  stock_quantity: number
  is_available: boolean
  contact_type: ContactType
  phone_number: string
  whatsapp_number: string
  status: ProductStatus
  rejection_note: string
  is_featured: boolean
  views_count: number
  phone_clicks: number
  whatsapp_clicks: number
  created_at: string
  updated_at: string
  seller?: Seller | null
  category?: Category | null
  images?: ProductImage[]
}

export interface Review {
  id: string
  product_id: string
  user_id: string
  rating: number
  comment: string
  status: ReviewStatus
  created_at: string
  updated_at: string
  user?: User | null
  product?: Product | null
}

export interface Favorite {
  id: string
  user_id: string
  product_id: string
  created_at: string
  product?: Product | null
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

export interface SuccessResponse<T> {
  data: T
}

export interface AdminTotals {
  total_products: number
  total_approved_products: number
  total_pending_products: number
  total_sellers: number
  total_users: number
  total_views: number
  total_phone_clicks: number
  total_whatsapp_clicks: number
  total_reviews: number
}

export interface SellerTotals {
  total_products: number
  approved_products: number
  pending_products: number
  total_views: number
  total_phone_clicks: number
  total_whatsapp_clicks: number
}

export type SettingsMap = Record<string, string>
