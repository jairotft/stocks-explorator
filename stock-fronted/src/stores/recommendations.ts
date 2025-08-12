import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Stock } from './stocks'


export interface Recommendation {
    stocks: Stock[]
    message: string
    
}

export const useRecommendationsStore = defineStore('recommendations', () => {

    const recommendation = ref<Recommendation | null>(null);
    const isLoading = ref<boolean>(false)
    const error = ref<string | null>(null)
    
    
    async function fetchRecommendations() {
        isLoading.value = true
        error.value = null
        
        try {
            const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/v1/api'
            const url = `${apiBaseUrl}/stocks/recommendations`
            console.log('Fetching recommendations from:', url)
            
            const response = await fetch(url)
            console.log('Response status:', response.status, response.statusText)
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: Failed to fetch recommendations`)
            }
            
            const res = await response.json();
            console.log('Full API response:', res)
            console.log('Response data:', res.data)
            
            // Verificar si la respuesta tiene la estructura esperada
            if (!res.data) {
                console.warn('No data property in response, using full response as data')
                recommendation.value = res as Recommendation
            } else {
                recommendation.value = res.data as Recommendation
            }
            
            console.log('Recommendation set to:', recommendation.value)
            return res.data || res
        } catch (err: any) {
            console.error('Error in fetchRecommendations:', err)
            error.value = err?.message ?? 'Error al cargar recomendaciones'
            throw err
        } finally {
            isLoading.value = false
        }
    }

    return {
        fetchRecommendations,
        recommendation,
        isLoading,
        error
    }
})