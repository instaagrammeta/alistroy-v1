import type { ProductFormState } from '~/components/ProductForm.vue'
import type { Product } from '~/types/api'

export function emptyProductForm(): ProductFormState {
  return {
    category_id: '', subcategory_id: '', brand_id: '',
    name_tj: '', name_ru: '', description_tj: '', description_ru: '',
    sku: '', unit: 'pcs', currency: 'TJS',
    sale_price: 0, cost_price: 0, discount_percent: 0,
    stock_quantity: 0, minimum_stock: 0, is_available: true, images: [],
  }
}

export function productToForm(p: Product): ProductFormState {
  return {
    category_id: p.category_id || '',
    subcategory_id: p.subcategory_id || '',
    brand_id: p.brand_id || '',
    name_tj: p.name_tj, name_ru: p.name_ru,
    description_tj: p.description_tj, description_ru: p.description_ru,
    sku: p.sku, unit: p.unit || 'pcs', currency: p.currency || 'TJS',
    sale_price: p.sale_price, cost_price: p.cost_price, discount_percent: p.discount_percent,
    stock_quantity: p.stock_quantity, minimum_stock: p.minimum_stock,
    is_available: p.is_available,
    images: (p.images || []).map((i) => ({ url: i.url, alt: i.alt })),
  }
}

export function formToPayload(f: ProductFormState): Record<string, any> {
  const payload: Record<string, any> = {
    category_id: f.category_id,
    name_tj: f.name_tj, name_ru: f.name_ru,
    description_tj: f.description_tj, description_ru: f.description_ru,
    sku: f.sku, unit: f.unit, currency: f.currency,
    sale_price: f.sale_price, cost_price: f.cost_price, discount_percent: f.discount_percent,
    stock_quantity: f.stock_quantity, minimum_stock: f.minimum_stock,
    is_available: f.is_available,
    images: f.images,
  }
  if (f.subcategory_id) payload.subcategory_id = f.subcategory_id
  if (f.brand_id) payload.brand_id = f.brand_id
  return payload
}
