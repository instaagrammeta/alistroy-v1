// Load public marketplace settings on app init (server + client).
export default defineNuxtPlugin(async () => {
  const settings = useSettingsStore()
  await settings.load()
})
