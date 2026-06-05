<script setup lang="ts">
import type { User } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const api = useApi()
const { t } = useI18n()
const id = route.params.id as string

const { data, refresh } = await useAsyncData(`admin-user-${id}`, () =>
  api<{ data: User }>(`/admin/users/${id}`)
)
const user = (data.value as any)?.data as User

const form = reactive({
  name: user.name,
  phone: user.phone,
  role: user.role,
  is_active: user.is_active,
  password: '',
})

const message = ref('')
const error = ref('')
const saving = ref(false)

const save = async () => {
  message.value = ''
  error.value = ''
  saving.value = true
  try {
    const body: Record<string, any> = {
      name: form.name,
      phone: form.phone,
      role: form.role,
      is_active: form.is_active,
    }
    if (form.password) body.password = form.password
    await api(`/admin/users/${id}`, { method: 'PATCH', body })
    form.password = ''
    message.value = t('common.save')
    await refresh()
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${user.email} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ user.email }}</h1>

    <div class="card p-6 max-w-xl">
      <form class="space-y-4" @submit.prevent="save">
        <div>
          <label class="label">{{ t('auth.email') }}</label>
          <input :value="user.email" disabled class="input bg-gray-50" />
        </div>
        <div>
          <label class="label">{{ t('auth.name') }}</label>
          <input v-model="form.name" class="input" />
        </div>
        <div>
          <label class="label">{{ t('auth.phone') }}</label>
          <input v-model="form.phone" class="input" />
        </div>
        <div>
          <label class="label">Role</label>
          <select v-model="form.role" class="select">
            <option value="customer">customer</option>
            <option value="seller">seller</option>
            <option value="admin">admin</option>
          </select>
        </div>
        <label class="inline-flex items-center gap-2">
          <input v-model="form.is_active" type="checkbox" class="h-4 w-4" />
          <span class="text-sm">Active</span>
        </label>
        <div>
          <label class="label">{{ t('auth.new_password') }} ({{ t('common.optional') }})</label>
          <input v-model="form.password" type="password" class="input" minlength="8" />
        </div>
        <div v-if="message" class="text-sm text-emerald-600">{{ message }}</div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? t('common.loading') : t('common.save') }}</button>
      </form>
    </div>
  </div>
</template>
