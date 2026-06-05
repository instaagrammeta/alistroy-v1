/**
 * Returns the localized version of a bilingual (TJ/RU) field.
 * Falls back to the other locale if one side is empty.
 */
export const useLocaleField = () => {
  const { locale } = useI18n()

  const pick = (tj?: string | null, ru?: string | null): string => {
    const t = (tj || '').trim()
    const r = (ru || '').trim()
    if (locale.value === 'ru') return r || t
    return t || r
  }

  const productTitle = (p?: { title_tj?: string; title_ru?: string } | null) =>
    pick(p?.title_tj, p?.title_ru)

  const productDescription = (p?: { description_tj?: string; description_ru?: string } | null) =>
    pick(p?.description_tj, p?.description_ru)

  const categoryTitle = (c?: { title_tj?: string; title_ru?: string } | null) =>
    pick(c?.title_tj, c?.title_ru)

  const sellerDescription = (s?: { description_tj?: string; description_ru?: string } | null) =>
    pick(s?.description_tj, s?.description_ru)

  return { pick, productTitle, productDescription, categoryTitle, sellerDescription }
}
