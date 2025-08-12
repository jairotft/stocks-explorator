import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Stock {
    code: string
    ticker: string
    company: string
    brokerage: string
    action: string
    rating_from: number
    rating_to: number
    target_from: number
    target_to: number
    crated_at: number
    updated_at: number
}

export interface PageMeta {
    page: number
    pageSize: number
    totalItems: number
    totalPages: number
    
}

export interface StockFilters {
    ticker?: string
    company?: string
    brokerage?: string
}

export const useStocksStore = defineStore('stocks', () => {
    const itemsById = ref<Record<string, Stock>>({})
    const pageIndex = ref<Record<number, string[]>>({})
    const loadedPages = ref<Set<number>>(new Set())

    const currentPage = ref<number>(1)
    const pageSize = ref<number>(5)

    const totalPages = ref<number>(1)
    const totalItems = ref<number>(0)

    const isLoading = ref<boolean>(false)
    const error = ref<string | null>(null)

    // Estados para los filtros
    const filters = ref<StockFilters>({})

    const currentList = computed<Stock[]>(() => {
        const page = pageIndex.value[currentPage.value] ?? []
        return page.map((id) => itemsById.value[id]).filter(Boolean)
    })

    const hasNextPage = computed<boolean>(() => {
        return currentPage.value < totalPages.value
    })

    const hasPreviousPage = computed<boolean>(() => {
        return currentPage.value > 1
    })

    function upsertMany(stocks: Stock[]) {
        stocks.forEach((stock) => {
            itemsById.value[stock.code] = stock
        })
    }

    function setPage(page: number, stocks: Stock[]) {
        pageIndex.value[page] = stocks.map(s => s.code)
        loadedPages.value.add(page)
    }

    async function fetchPageFromApi(page: number, currentFilters?: StockFilters) {
        // Construir la URL base usando variable de entorno
        const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/v1/api'
        let url = `${apiBaseUrl}/stocks/list?page=${page}&pageSize=${pageSize.value}`
        
        // Agregar filtros como query parameters solo si tienen valores
        const filtersToUse = currentFilters || filters.value
        
        if (filtersToUse.ticker && filtersToUse.ticker.trim()) {
            url += `&ticker=${encodeURIComponent(filtersToUse.ticker.trim())}`
        }
        
        if (filtersToUse.company && filtersToUse.company.trim()) {
            url += `&company=${encodeURIComponent(filtersToUse.company.trim())}`
        }
        
        if (filtersToUse.brokerage && filtersToUse.brokerage.trim()) {
            url += `&brokerage=${encodeURIComponent(filtersToUse.brokerage.trim())}`
        }

        const response = await fetch(url)
        if (!response.ok) {
            throw new Error('Failed to fetch stocks')
        }
        const res = await response.json();
        console.log('fetchPageFromApi', page, 'data', res.data, 'filters', filtersToUse)
        const data = res.data

        const body: {
            items: Stock[];
            meta: PageMeta;
        } = {
            items: data.stocks ?? [],
            meta: {
                page: data.current_page,
                pageSize: data.per_page,
                totalItems: data.total,
                totalPages: Math.ceil(data.total / data.per_page)
            }
        }

        return body
    }

    async function ensurePage(page: number, {force = false} = {}) {
        
        if (!force && loadedPages.value.has(page)) {
            console.log('ensurePage', page, 'already loaded')
            return
        }
        console.log('ensurePage', page, 'loading')

        isLoading.value = true
        error.value = null

        try {
            const {items, meta} = await fetchPageFromApi(page)
            console.log('ensurePage', page, 'items', items)
            upsertMany(items)
            setPage(page, items)
            totalPages.value = meta.totalPages
            totalItems.value = meta.totalItems
        } catch (err: any   ) {
            error.value = err?.message ?? 'Error al cargar stocks'
            throw err
        } finally {
            isLoading.value = false
        }
    }

    async function goToPage(page: number) {
        if (page < 1 || page > totalPages.value) {
            throw new Error('Page out of bounds')
        }

        currentPage.value = page
        await ensurePage(page)
    }

    async function nextPage() {
        if (hasNextPage.value) {
            await goToPage(currentPage.value + 1)
        }
    }

    async function previousPage() {
        if (hasPreviousPage.value) {
            await goToPage(currentPage.value - 1)
        }
    }

    async function prefetchNextPage() {
        if (hasNextPage.value &&  !loadedPages.value.has(currentPage.value + 1)) {
            try {
                await ensurePage(currentPage.value + 1)
            } catch (err) {
                console.error('Error prefetching next page', err)
            }
        }
    }

    // Función para actualizar filtros
    function updateFilters(newFilters: StockFilters) {
        filters.value = { ...newFilters }
    }

    // Función para aplicar filtros y recargar datos
    async function applyFilters(newFilters: StockFilters) {
        updateFilters(newFilters)
        clearAll()
        await goToPage(1)
    }

    // Función para limpiar filtros
    async function clearFilters() {
        filters.value = {}
        clearAll()
        await goToPage(1)
    }

    function clearAll() {
        itemsById.value = {}
        pageIndex.value = {}
        loadedPages.value = new Set()
        totalPages.value = 1
        totalItems.value = 0
        currentPage.value = 1
        error.value = null
        
    }
    

    return {
        itemsById,
        pageIndex,
        loadedPages,
        currentPage,
        pageSize,
        totalPages,
        totalItems,
        isLoading,
        error,
        filters,
        currentList,
        hasNextPage,
        hasPreviousPage,
        ensurePage,
        goToPage,
        nextPage,
        previousPage,
        prefetchNextPage,
        updateFilters,
        applyFilters,
        clearFilters,
        clearAll,
    }  
    
})