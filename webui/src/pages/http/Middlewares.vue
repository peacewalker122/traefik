<template>
  <page-default>
    <section class="app-section">
      <div class="app-section-wrap app-boxed app-boxed-xl q-pl-md q-pr-md q-pt-xl q-pb-xl">
        <div class="row no-wrap items-center q-mb-lg">
          <tool-bar-table
            v-model:status="status"
            v-model:filter="filter"
          />
        </div>
        <div class="row items-center q-col-gutter-lg">
          <div class="col-12">
            <main-table
              ref="mainTable"
              v-bind="getTableProps({ type: 'http-middlewares' })"
              :data="allMiddlewares.items"
              :on-load-more="handleLoadMore"
              :end-reached="allMiddlewares.endReached"
              :loading="allMiddlewares.loading"
            />
          </div>
        </div>
      </div>
    </section>
  </page-default>
</template>

<script>
import { defineComponent } from 'vue'
import { mapActions, mapGetters } from 'vuex'
import GetTablePropsMixin from '../../_mixins/GetTableProps'
import PaginationMixin from '../../_mixins/Pagination'
import PageDefault from '../../components/_commons/PageDefault.vue'
import ToolBarTable from '../../components/_commons/ToolBarTable.vue'
import MainTable from '../../components/_commons/MainTable.vue'

export default defineComponent({
  name: 'PageHTTPMiddlewares',
  components: {
    PageDefault,
    ToolBarTable,
    MainTable
  },
  mixins: [
    GetTablePropsMixin,
    PaginationMixin({
      fetchMethod: 'getAllMiddlewaresWithParams',
      scrollerRef: 'mainTable.$refs.scroller',
      pollingIntervalTime: 5000
    })
  ],
  data () {
    return {
      filter: '',
      status: ''
    }
  },
  computed: {
    ...mapGetters('http', { allMiddlewares: 'allMiddlewares' })
  },
  watch: {
    'status' () {
      this.refreshAll()
    },
    'filter' () {
      this.refreshAll()
    }
  },
  beforeUnmount () {
    this.$store.commit('http/getAllMiddlewaresClear')
  },
  methods: {
    ...mapActions('http', { getAllMiddlewares: 'getAllMiddlewares' }),
    getAllMiddlewaresWithParams (params) {
      return this.getAllMiddlewares({
        query: this.filter,
        status: this.status,
        ...params
      })
    },
    refreshAll () {
      if (this.allMiddlewares.loading) {
        return
      }

      this.initFetch()
    },
    handleLoadMore ({ page = 1 } = {}) {
      return this.fetchMore({ page })
    }
  }
})
</script>

<style scoped lang="scss">

</style>
