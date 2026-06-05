// Loads public marketplace settings (logo, names, hero, etc.) on app init.
export default defineNuxtPlugin(async () => {
  const settings = useSettingsStore()
  await settings.load()
})
