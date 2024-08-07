<template>
  <page-default>
    <section class="app-section">
      <div class="app-section-wrap app-boxed app-boxed-xl q-pl-md q-pr-md q-pt-xl q-pb-xl">
        <div
          v-if="!loading"
          class="row items-start"
        >
          <div
            v-if="entryPoints.length"
            class="col-12 col-md-3 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-log-in-outline" />
              <div class="app-title-label">
                Entrypoints
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12 col-md-8">
                <div class="row items-start q-col-gutter-md">
                  <div
                    v-for="(entryPoint, index) in entryPoints"
                    :key="index"
                    class="col-12"
                  >
                    <panel-entry
                      type="detail"
                      ex-size="true"
                      :name="entryPoint.name"
                      :address="entryPoint.address"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12 col-md-4 xs-hide sm-hide">
                <q-icon
                  name="eva-arrow-forward-outline"
                  class="arrow"
                />
              </div>
            </div>
          </div>

          <div
            v-if="routerByName.item.name"
            class="col-12 col-md-3 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-globe-outline" />
              <div class="app-title-label">
                {{ routerType }}
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12 col-md-8">
                <div class="row items-start q-col-gutter-md">
                  <div class="col-12">
                    <panel-entry
                      focus="true"
                      type="detail"
                      name="router"
                      :address="routerByName.item.name"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12 col-md-4 xs-hide sm-hide">
                <q-icon
                  name="eva-arrow-forward-outline"
                  class="arrow"
                />
              </div>
            </div>
          </div>

          <div
            v-if="hasMiddlewares"
            class="col-12 col-md-3 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-layers" />
              <div class="app-title-label">
                {{ middlewareType }}
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12 col-md-8">
                <div class="row items-start q-col-gutter-md">
                  <div
                    v-for="(middleware, index) in middlewares"
                    :key="index"
                    class="col-12"
                  >
                    <panel-entry
                      type="detail"
                      name="Middleware"
                      :address="middleware.type"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12 col-md-4 xs-hide sm-hide">
                <q-icon
                  name="eva-arrow-forward-outline"
                  class="arrow"
                />
              </div>
            </div>
          </div>

          <div
            v-if="routerByName.item.service"
            class="service col-12 col-md-3 q-mb-lg path-block"
            @click="$router.push({ path: `/${protocol}/services/${getServiceId(routerByName.item)}`})"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-flash" />
              <div class="app-title-label">
                Service
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12 col-md-8">
                <div class="row items-start q-col-gutter-md">
                  <div class="col-12">
                    <panel-entry
                      type="detail"
                      name="Service"
                      :address="routerByName.item.service"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          v-else
          class="row items-start"
        >
          <div class="col-12">
            <p
              v-for="n in 4"
              :key="n"
              class="flex"
            >
              <SkeletonBox
                :min-width="15"
                :max-width="15"
                style="margin-right: 2%"
              /> <SkeletonBox
                :min-width="50"
                :max-width="83"
              />
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="app-section">
      <div class="app-section-wrap app-boxed app-boxed-xl q-pl-md q-pr-md q-pt-xl q-pb-xl">
        <div
          v-if="!loading"
          class="row items-start q-col-gutter-md"
        >
          <div
            v-if="routerByName.item"
            class="col-12 col-md-4 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-info" />
              <div class="app-title-label">
                Router Details
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12">
                <div class="row items-start q-col-gutter-md">
                  <div class="col-12">
                    <panel-router-details
                      :data="routerByName.item"
                      :protocol="protocol"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="protocol !== 'udp'"
            class="col-12 col-md-4 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-shield" />
              <div class="app-title-label">
                TLS
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12">
                <div class="row items-start q-col-gutter-md">
                  <div class="col-12">
                    <panel-t-l-s
                      :data="routerByName.item.tls"
                      :protocol="protocol"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="protocol !== 'udp'"
            class="col-12 col-md-4 q-mb-lg path-block"
          >
            <div class="row no-wrap items-center q-mb-lg app-title">
              <q-icon name="eva-layers" />
              <div class="app-title-label">
                Middlewares
              </div>
            </div>
            <div class="row items-start q-col-gutter-lg">
              <div class="col-12">
                <div class="row items-start q-col-gutter-md">
                  <div class="col-12">
                    <panel-middlewares :data="middlewares" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          v-else
          class="row items-start"
        >
          <div class="col-12">
            <p
              v-for="n in 4"
              :key="n"
              class="flex"
            >
              <SkeletonBox
                :min-width="15"
                :max-width="15"
                style="margin-right: 2%"
              /> <SkeletonBox
                :min-width="50"
                :max-width="83"
              />
            </p>
          </div>
        </div>
      </div>
    </section>
  </page-default>
</template>

<script>
import { defineComponent } from 'vue'
import { mapActions, mapGetters } from 'vuex'
import PageDefault from '../../components/_commons/PageDefault.vue'
import SkeletonBox from '../../components/_commons/SkeletonBox.vue'
import PanelEntry from '../../components/dashboard/PanelEntry.vue'
import PanelRouterDetails from '../../components/_commons/PanelRouterDetails.vue'
import PanelTLS from '../../components/_commons/PanelTLS.vue'
import PanelMiddlewares from '../../components/_commons/PanelMiddlewares.vue'

export default defineComponent({
  name: 'PageRouterDetail',
  components: {
    PageDefault,
    SkeletonBox,
    PanelEntry,
    PanelRouterDetails,
    PanelTLS,
    PanelMiddlewares
  },
  props: {
    name: { type: String, default: undefined, required: false },
    type: { type: String, default: undefined, required: false }
  },
  data () {
    return {
      loading: true,
      entryPoints: [],
      middlewares: [],
      timeOutGetAll: null
    }
  },
  computed: {
    hasTLSConfiguration () {
      return this.routerByName.item.tls
    },
    middlewareType () {
      return this.$route.meta.protocol.toUpperCase() + ' Middlewares'
    },
    routerType () {
      return this.$route.meta.protocol.toUpperCase() + ' Router'
    },
    ...mapGetters('http', { http_routerByName: 'routerByName' }),
    ...mapGetters('tcp', { tcp_routerByName: 'routerByName' }),
    ...mapGetters('udp', { udp_routerByName: 'routerByName' }),
    hasMiddlewares () {
      return this.$route.meta.protocol !== 'udp' && this.middlewares.length > 0
    },
    protocol () {
      return this.$route.meta.protocol
    },
    routerByName () {
      return this[`${this.protocol}_routerByName`]
    },
    getRouterByName () {
      return this[`${this.protocol}_getRouterByName`]
    },
    getMiddlewareByName () {
      return this[`${this.protocol}_getMiddlewareByName`]
    }
  },
  created () {
    this.refreshAll()
  },
  beforeUnmount () {
    clearInterval(this.timeOutGetAll)
    this.$store.commit('http/getRouterByNameClear')
    this.$store.commit('tcp/getRouterByNameClear')
    this.$store.commit('udp/getRouterByNameClear')
  },
  methods: {
    ...mapActions('http', { http_getRouterByName: 'getRouterByName', http_getMiddlewareByName: 'getMiddlewareByName' }),
    ...mapActions('tcp', { tcp_getRouterByName: 'getRouterByName', tcp_getMiddlewareByName: 'getMiddlewareByName' }),
    ...mapActions('udp', { udp_getRouterByName: 'getRouterByName' }),
    ...mapActions('entrypoints', { getEntrypointsByName: 'getByName' }),
    refreshAll () {
      if (this.routerByName.loading) {
        return
      }
      this.onGetAll()
    },
    onGetAll () {
      this.getRouterByName(this.name)
        .then(body => {
          if (!body) {
            this.loading = false
            return
          }
          // Get entryPoints
          if (body.using) {
            for (const entryPoint in body.using) {
              if (Object.getOwnPropertyDescriptor(body.using, entryPoint)) {
                this.getEntrypointsByName(body.using[entryPoint])
                  .then(body => {
                    if (body) {
                      this.entryPoints.push(body)
                    }
                  })
                  .catch(error => {
                    console.log('Error -> entrypoints/byName', error)
                  })
              }
            }
          }
          // Get middlewares
          if (body.middlewares) {
            for (const middleware in body.middlewares) {
              if (Object.getOwnPropertyDescriptor(body.middlewares, middleware)) {
                this.getMiddlewareByName(body.middlewares[middleware])
                  .then(body => {
                    if (body) {
                      this.middlewares.push(body)
                    }
                  })
                  .catch(error => {
                    console.log('Error -> middlewares/byName', error)
                  })
              }
            }
          }
          clearTimeout(this.timeOutGetAll)
          this.timeOutGetAll = setTimeout(() => {
            this.loading = false
          }, 300)
        })
        .catch(error => {
          console.log('Error -> routers/byName', error)
        })
    },
    getServiceId (data) {
      const words = data.service.split('@')
      if (words.length === 2) {
        return data.service
      }

      return `${encodeURIComponent(data.service)}@${data.provider}`
    }
  }
})
</script>

<style scoped lang="scss">
  @import "../../css/sass/variables";

  .path-block {
    .arrow {
      font-size: 40px;
      margin-top: 20px;
      margin-left: 20px;
      color: #b2b2b2;
    }

    &.service {
      cursor: pointer;
    }
  }

</style>
