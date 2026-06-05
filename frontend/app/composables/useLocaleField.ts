/** Picks the localized variant of a bilingual TJ/RU field with fallback. */
export const useLocaleField = () => {
  const { locale } = useI18n()

  const pick = (tj?: string | null, ru?: string | null): string => {
    const t = (tj || '').trim()
    const r = (ru || '').trim()
    if (locale.value === 'ru') return r || t
    return t || r
  }

  const productName = (p?: { name_tj?: string; name_ru?: string } | null) => pick(p?.name_tj, p?.name_ru)
  const productDesc = (p?: { description_tj?: string; description_ru?: string } | null) => pick(p?.description_tj, p?.description_ru)
  const categoryName = (c?: { name_tj?: string; name_ru?: string } | null) => pick(c?.name_tj, c?.name_ru)

  return { pick, productName, productDesc, categoryName }
}
