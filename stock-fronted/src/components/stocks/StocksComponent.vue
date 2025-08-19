<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useStocksStore } from '../../stores/stocks'
import { useRouter } from 'vue-router'

const stocks = useStocksStore()
const router = useRouter()

// Estados para los filtros (locales para el input)
const filterTicker = ref('')
const filterCompany = ref('')
const filterBrokerage = ref('')

// Estados para el ordenamiento
const currentSortField = ref<string>('')
const currentSortDirection = ref<string>('asc')

// Función helper para extraer valores de objetos de validación
const extractValue = (value: any): string => {
  if (typeof value === 'object' && value !== null && 'String' in value) {
    return value.String || ''
  }
  return value || ''
}

// Función helper para formatear valores monetarios de manera segura
const formatCurrency = (value: number | string | null | undefined): string => {
  if (value === null || value === undefined || value === '') {
    return 'N/A'
  }
  
  const numValue = typeof value === 'string' ? parseFloat(value) : value
  
  if (isNaN(numValue)) {
    return 'N/A'
  }
  
  return `$${numValue.toFixed(2)}`
}

/*
  * Función para formatear fechas
  * De un formato "2025-08-07T19:30:09.885518-05:00"
  * Convertimos a una fecha legible dd/mm/yyyy hh:mm:ss
  */
const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  
  if (isNaN(date.getTime())) {
    return 'N/A'
  }
  
  const options: Intl.DateTimeFormatOptions = {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZoneName: 'short'
  }
  
  return date.toLocaleString('es-ES', options)
  
}

// Función para aplicar filtros con debounce
let filterTimeout: ReturnType<typeof setTimeout> | null = null

const applyFiltersWithDebounce = () => {
  if (filterTimeout) {
    clearTimeout(filterTimeout)
  }
  
  filterTimeout = setTimeout(async () => {
    const filters = {
      ticker: filterTicker.value,
      company: filterCompany.value,
      brokerage: filterBrokerage.value
    }
    
    await stocks.applyFilters(filters)
  }, 500) // Debounce de 500ms
}

// Watchers para aplicar filtros automáticamente
watch([filterTicker, filterCompany, filterBrokerage], applyFiltersWithDebounce)

// Función para limpiar todos los filtros
const clearFilters = async () => {
  filterTicker.value = ''
  filterCompany.value = ''
  filterBrokerage.value = ''
  await stocks.clearFilters()
}

const goToRecommendationsPage = () => {
  router.push('/recommendations')
}

// Función para manejar el ordenamiento
const handleSort = async (field: string) => {
  let direction = '0'
  
  // Si ya estamos ordenando por este campo, cambiar la dirección
  if (currentSortField.value === field) {
    direction = currentSortDirection.value === '0' ? '1' : '0'
  }
  
  currentSortField.value = field
  currentSortDirection.value = direction

  // Aplicar ordenamiento usando el store
  const orderDirection = direction === '0' ? '1' : '0'
  await stocks.applyOrdering({
    orderBy: field,
    orderDirection: orderDirection
  })
}

// Función para obtener el icono de ordenamiento
const getSortIcon = (field: string) => {
  if (currentSortField.value !== field) {
    return 'sort'
  }
  return currentSortDirection.value === 'asc' ? 'sort-up' : 'sort-down'
}

// Inicializar filtros desde el store si existen
onMounted(() => {
  filterTicker.value = stocks.filters.ticker || ''
  filterCompany.value = stocks.filters.company || ''
  filterBrokerage.value = stocks.filters.brokerage || ''
  
  stocks.goToPage(1).then(() => {
    stocks.prefetchNextPage()
  })
})
</script>


<template>
  <div class="p-6 max-w-7xl mx-auto">
    <div class="bg-white shadow-lg rounded-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <h1 class="text-2xl font-bold text-gray-900">
          Stocks (página {{ stocks.currentPage }} / {{ stocks.totalPages }})
        </h1>
      </div>

      <!-- Sección de filtros -->
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-25">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
          <div>
            <label for="filter-ticker" class="block text-sm font-medium text-gray-700 mb-1">
              Filtrar por Ticker
            </label>
            <input
              id="filter-ticker"
              v-model="filterTicker"
              type="text"
              placeholder="Ej: AKBA"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
            />
          </div>
          <div>
            <label for="filter-company" class="block text-sm font-medium text-gray-700 mb-1">
              Filtrar por Compañia
            </label>
            <input
              id="filter-company"
              v-model="filterCompany"
              type="text"
              placeholder="Ej: Akebia"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
            />
          </div>
          <div>
            <label for="filter-brokerage" class="block text-sm font-medium text-gray-700 mb-1">
              Filtrar por Broker
            </label>
            <input
              id="filter-brokerage"
              v-model="filterBrokerage"
              type="text"
              placeholder="Ej: HC Wainwright"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
            />
          </div>
          <div>
            <button
              @click="clearFilters"
              class="w-full inline-flex justify-center items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
              </svg>
              Limpiar filtros
            </button>
          </div>
        </div>
        
        <!-- Contador de resultados -->
        <div v-if="filterTicker || filterCompany || filterBrokerage" class="mt-3 text-sm text-gray-600">
          Total de {{ stocks.totalItems }} resultados encontrados
        </div>
      </div>

      <div v-if="stocks.isLoading" class="p-8 text-center">
        <div class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-lg text-gray-600">Cargando…</span>
        </div>
      </div>

      <div v-else-if="stocks.error" class="p-8 text-center">
        <div class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <div class="mt-2 text-sm text-red-700">{{ stocks.error }}</div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="stocks.currentList.length === 0" class="p-8 text-center">
        <div class="text-gray-500">
          <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
            <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No se encontraron resultados</h3>
          <p class="mt-1 text-sm text-gray-500">Intenta ajustar los filtros para ver más resultados.</p>
        </div>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th 
                @click="handleSort('ticker')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors duration-200"
              >
                <div class="flex items-center space-x-1">
                  <span>Ticker</span>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path 
                      v-if="getSortIcon('ticker') === 'sort'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                    />
                    <path 
                      v-else-if="getSortIcon('ticker') === 'sort-up'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M5 15l7-7 7 7"
                    />
                    <path 
                      v-else
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </div>
              </th>
              <th 
                @click="handleSort('company')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors duration-200"
              >
                <div class="flex items-center space-x-1">
                  <span>Company</span>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path 
                      v-if="getSortIcon('company') === 'sort'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                    />
                    <path 
                      v-else-if="getSortIcon('company') === 'sort-up'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M5 15l7-7 7 7"
                    />
                    <path 
                      v-else
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </div>
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Brokerage
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Action
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Rating From
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Rating To
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Target From
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Target To
              </th>
              <th 
                @click="handleSort('record_time')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors duration-200"
              >
                <div class="flex items-center space-x-1">
                  <span>Record Time</span>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path 
                      v-if="getSortIcon('record_time') === 'sort'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                    />
                    <path 
                      v-else-if="getSortIcon('record_time') === 'sort-up'"
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M5 15l7-7 7 7"
                    />
                    <path 
                      v-else
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </div>
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="item in stocks.currentList" :key="item.code" class="hover:bg-gray-50 transition-colors duration-200">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {{ extractValue(item.ticker) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ extractValue(item.company) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                {{ extractValue(item.brokerage) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                  {{ extractValue(item.action) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                  {{ extractValue(item.rating_from) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                  {{ extractValue(item.rating_to) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600">
                {{ formatCurrency(item.target_from) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600">
                {{ formatCurrency(item.target_to) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(item.record_time) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginación -->
      <div class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
        <div class="flex items-center justify-between">
          <div class="flex-1 flex justify-between sm:hidden">
            <button
              :disabled="!stocks.hasPreviousPage"
              @click="stocks.previousPage()"
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Anterior
            </button>
            <button
              :disabled="!stocks.hasNextPage"
              @click="stocks.nextPage()"
              class="relative ml-3 inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Siguiente
            </button>
          </div>
          <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700">
                Página
                <span class="font-medium">{{ stocks.currentPage }}</span>
                de
                <span class="font-medium">{{ stocks.totalPages }}</span>
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                <button
                  :disabled="!stocks.hasPreviousPage"
                  @click="stocks.previousPage()"
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span class="sr-only">Anterior</span>
                  <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button
                  :disabled="!stocks.hasNextPage"
                  @click="stocks.nextPage()"
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span class="sr-only">Siguiente</span>
                  <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
      <!-- Boton para ir a pagina de recomendaciones-->
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-25">
        <button
          @click="goToRecommendationsPage()"
          class="w-full inline-flex justify-center items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          Ver recomendaciones
        </button>
      </div>
    </div>
  </div>
</template>
