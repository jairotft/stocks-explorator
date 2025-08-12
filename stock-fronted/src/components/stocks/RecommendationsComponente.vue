<script setup lang="ts">
import { useRecommendationsStore } from '../../stores/recommendations'
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'

const recommendationsStore = useRecommendationsStore()
const router = useRouter()

// Configurar markdown-it
const md = new MarkdownIt({
  html: true,        // Permitir HTML en el markdown
  linkify: true,     // Convertir URLs automáticamente en enlaces
  typographer: true, // Reemplazar comillas y guiones por versiones tipográficas
  breaks: true       // Convertir saltos de línea en <br>
})

// Función helper para extraer valores de objetos de validación
const extractValue = (value: any): string => {
  if (typeof value === 'object' && value !== null && 'String' in value) {
    return value.String || ''
  }
  return value || ''
}

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

// Computed para el mensaje formateado usando markdown-it
const formattedMessage = computed(() => {
  const message = recommendationsStore.recommendation?.message
  if (!message) return ''
  
  // Renderizar markdown a HTML
  const html = md.render(message)
  
  // Agregar clases de Tailwind a los elementos HTML generados
  return html
    .replace(/<h1>/g, '<h1 class="text-2xl font-bold text-gray-900 mb-4">')
    .replace(/<h2>/g, '<h2 class="text-xl font-semibold text-gray-900 mb-3">')
    .replace(/<h3>/g, '<h3 class="text-lg font-semibold text-gray-900 mb-2">')
    .replace(/<p>/g, '<p class="text-gray-700 leading-relaxed mb-4">')
    .replace(/<strong>/g, '<strong class="font-semibold text-gray-900">')
    .replace(/<em>/g, '<em class="italic text-gray-700">')
    .replace(/<a /g, '<a class="text-blue-600 hover:text-blue-800 underline" target="_blank" rel="noopener noreferrer" ')
    .replace(/<code>/g, '<code class="bg-gray-100 text-gray-800 px-1 py-0.5 rounded text-sm font-mono">')
    .replace(/<ul>/g, '<ul class="list-disc list-inside mb-4 space-y-1">')
    .replace(/<ol>/g, '<ol class="list-decimal list-inside mb-4 space-y-1">')
    .replace(/<li>/g, '<li class="text-gray-700">')
    .replace(/<blockquote>/g, '<blockquote class="border-l-4 border-blue-200 pl-4 py-2 mb-4 bg-blue-50 italic text-gray-700">')
    .replace(/<pre>/g, '<pre class="bg-gray-100 rounded-lg p-4 mb-4 overflow-x-auto">')
    .replace(/<table>/g, '<table class="min-w-full table-auto mb-4 border-collapse">')
    .replace(/<th>/g, '<th class="border border-gray-300 px-4 py-2 bg-gray-100 font-semibold text-left">')
    .replace(/<td>/g, '<td class="border border-gray-300 px-4 py-2 text-gray-700">')
})

const goBackToStocks = () => {
  router.push('/stocks')
}

onMounted(() => {
    recommendationsStore.fetchRecommendations().then(() => {
        console.log('Recommendations fetched successfully')
    }).catch((error) => {
        console.error('Error fetching recommendations:', error)
    })
})
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto">
    <div class="bg-white shadow-lg rounded-lg overflow-hidden">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">
            Recomendaciones de Inversión
          </h1>
          <div class="flex space-x-2">
            <button
              @click="goBackToStocks"
              class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
              Volver a Stocks
            </button>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="recommendationsStore.isLoading" class="p-8 text-center">
        <div class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-lg text-gray-600">Cargando recomendaciones…</span>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="recommendationsStore.error" class="p-8 text-center">
        <div class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <div class="mt-2 text-sm text-red-700">{{ recommendationsStore.error }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div v-else-if="recommendationsStore.recommendation">
        <!-- Mensaje de recomendación -->
        <div class="px-6 py-6 bg-gradient-to-r from-blue-50 to-indigo-50 border-b border-gray-200">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-blue-600 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="ml-3 flex-1">
              <h3 class="text-lg font-medium text-gray-900 mb-3">
                Análisis de Recomendaciones
              </h3>
              <div 
                class="prose prose-sm max-w-none markdown-content"
                v-html="formattedMessage"
              ></div>
            </div>
          </div>
        </div>

        <!-- Tabla de stocks recomendados -->
        <div v-if="recommendationsStore.recommendation.stocks && recommendationsStore.recommendation.stocks.length > 0" class="overflow-x-auto">
          <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-900">
              Stocks Recomendados ({{ recommendationsStore.recommendation.stocks.length }})
            </h3>
          </div>
          
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Ticker
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Company
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
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="stock in recommendationsStore.recommendation.stocks" :key="extractValue(stock.ticker)" class="hover:bg-gray-50 transition-colors duration-200">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {{ extractValue(stock.ticker) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ extractValue(stock.company) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                  {{ extractValue(stock.brokerage) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                  {{ extractValue(stock.action) }}
                </span>
              </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                    {{ extractValue(stock.rating_from) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                    {{ extractValue(stock.rating_to) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600">
                  {{ formatCurrency(stock.target_from) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600">
                  {{ formatCurrency(stock.target_to) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Estado vacío para stocks -->
        <div v-else class="p-8 text-center">
          <div class="text-gray-500">
            <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
              <path d="M34 40h10v-4a6 6 0 00-10.712-3.714M34 40H14m20 0v-4a9.971 9.971 0 00-.712-3.714M14 40H4v-4a6 6 0 0110.713-3.714M14 40v-4c0-1.313.253-2.566.713-3.714m0 0A10.003 10.003 0 0124 26c4.21 0 7.813 2.602 9.288 6.286M30 14a6 6 0 11-12 0 6 6 0 0112 0zm12 6a4 4 0 11-8 0 4 4 0 018 0zm-28 0a4 4 0 11-8 0 4 4 0 018 0z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No hay stocks recomendados</h3>
            <p class="mt-1 text-sm text-gray-500">En este momento no hay recomendaciones de inversión disponibles.</p>
          </div>
        </div>
      </div>

      <!-- Estado vacío general -->
      <div v-else class="p-8 text-center">
        <div class="text-gray-500">
          <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
            <path d="M8 14v20c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252M8 14c0 4.418 7.163 8 16 8s16-3.582 16-8M8 14c0-4.418 7.163-8 16-8s16 3.582 16 8m0 0v14m-16-5c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No hay datos disponibles</h3>
          <p class="mt-1 text-sm text-gray-500">No se pudieron cargar las recomendaciones en este momento.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Estilos específicos para contenido markdown */
.markdown-content {
  line-height: 1.6;
}

.markdown-content h1,
.markdown-content h2,
.markdown-content h3 {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}

.markdown-content h1:first-child,
.markdown-content h2:first-child,
.markdown-content h3:first-child {
  margin-top: 0;
}

.markdown-content p {
  margin-bottom: 1rem;
}

.markdown-content p:last-child {
  margin-bottom: 0;
}

.markdown-content ul,
.markdown-content ol {
  margin-bottom: 1rem;
  padding-left: 1.5rem;
}

.markdown-content li {
  margin-bottom: 0.25rem;
}

.markdown-content code {
  font-size: 0.875rem;
}

.markdown-content a:hover {
  text-decoration: underline;
}

.markdown-content strong {
  font-weight: 600;
}

.markdown-content em {
  font-style: italic;
}
</style>