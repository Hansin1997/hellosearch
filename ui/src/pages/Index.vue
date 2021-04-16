<template>
  <div>
    <q-page v-if="data.length==0" padding class="row items-center justify-evenly">
      <div style="max-width:600px" class="full-width">
        <div class="text-center text-h4 q-mb-md" style="margin-top: -100px">
          <q-icon name="search" class="q-mr-lg"/>
          <span>{{ $t("title") }}</span>
        </div>
        <q-form>
          <q-input v-model="q" outlined>
            <template v-slot:before>
              <q-select
                outlined
                v-model="i"
                :label="$t('index')"
                :options="indices"
                @filter="loadIndices"
                multiple
                style="width: 150px"
              >
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section class="text-grey">
                      {{ $t("noResults") }}
                    </q-item-section>
                  </q-item>
                </template>
              </q-select>
            </template>
            <template v-slot:append>
              <q-btn :loading="searching" type="submit" @click="search" round flat icon="search"/>
            </template>
          </q-input>
        </q-form>
      </div>
    </q-page>
    <q-page padding v-else>
      <div v-if="data.total==0" class="text-center q-mt-lg q-mb-lg text-grey">
        {{ $t("noResults") }}
      </div>
      <div v-else>
        <div class="q-pl-md text-left text-grey">
          <span>{{ (p - 1) * 10 }}-{{ Math.min(p * 10, data.total) }}</span>
          <span> / </span>
          <span>{{ data.total }}</span>
        </div>
        <!--        <q-separator class="q-mt-sm q-mb-sm"/>-->
        <q-list style="max-width: 800px" class="q-gutter-md">
          <q-item v-for="(item,index) in data.hits" :key="index">
            <q-item-section>
              <q-item-label class="text-h6 search-result">
                <a target="_blank" :href="buildHref(item)">
                  <span v-html="$xss(item.highlights.Title[0])"
                        v-if="item.highlights.Title && item.highlights.Title.length > 0"/>
                  <span v-else>
                  {{ item.Title }}
                </span>
                </a>
              </q-item-label>
              <q-item-label class="q-pa-none q-ma-none">
                <span class="link text-subtitle1">{{ item.Url }}</span>
              </q-item-label>
              <q-item-label>
                <div class="search-result content">
                  <span>{{ buildDate(item) }}</span>
                  <span class="q-ml-sm q-mr-sm">·</span>
                  <span v-html="$xss(buildContent(item))"/>
                </div>
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
        <div class="q-pa-lg" v-if="pageCount>1">
          <q-pagination
            :disable="searching"
            v-model="p"
            :max="pageCount"
            :max-pages="6"
            :direction-links="true"
          >

          </q-pagination>
        </div>
      </div>
    </q-page>
  </div>

</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  name: 'PageIndex',
  data() {
    return {
      data: [],
      indices: [],
      q: <string>this.$route.query.q,
      searching: false,
      p: (this.$route.query.p && Number(this.$route.query.p) > 0) ? Number(this.$route.query.p) : 1,
      i: (this.$route.query.i) ? [this.$route.query.i] : []
    }
  },
  methods: {
    refresh() {
      this.p = 1
      this.search()
    },
    search() {
      if (this.q != "" && this.q == this.$route.query.q && this.p == Number(this.$route.query.p) && this.buildQueryArray(this.i) == this.$route.query.i) {
        this.doSearch()
      } else
        this.$router.push({query: {q: this.q, p: "" + this.p, i: this.buildQueryArray(this.i)}})
    },
    doSearch() {
      if (this.searching) {
        // this.$root.$emit("finishSearch")
        return
      }
      this.searching = true
      this.$axios.get("/v1/doc", {
        params: {
          i: this.i && this.i.length > 0 ? this.buildQueryArray(this.i) : "spider-*",
          q: encodeURIComponent(this.q),
          p: this.p
        }
      }).then(res => {
        this.data = res.data
      }).finally(() => {
        this.searching = false
        this.$root.$emit("finishSearch")
      })
    },
    loadIndices(val: any, update: any, abort: any) {
      this.$axios.get("/v1/indices")
        .then(res => {
          update(() => {
            this.indices = res.data
          })
        }).catch((e) => {
        update(() => {
          this.indices = []
        })
        return Promise.reject(e)
      })
    },
    buildHref(item: any): string {
      return item.Url;
    },
    buildContent(item: any): string {
      let result = ""
      if (item.highlights.Content) {
        item.highlights.Content.forEach((c: string) => {
          result += c + " "
        })
      }
      return result
    },
    buildDate(item: any): string {
      let date: Date;
      if (item.ResponseHeader && item.ResponseHeader["Last-Modified"] && item.ResponseHeader["Last-Modified"].length)
        date = new Date(item.ResponseHeader["Last-Modified"][0])
      else
        date = new Date(item.FetchAt)
      return date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDate()
    },
    buildQueryArray(arr: Array<any> | null | undefined): string {
      let r = "";
      if (arr != null)
        arr.forEach(i => {
          if (r.length > 0)
            r += ","
          r += i
        })
      return r;
    },
    onQueryUpdate() {
      if (this.$route.query.q) {
        this.doSearch()
        this.$root.$emit("showSearch", this.q, true)
      } else {
        this.data = []
        this.p = 1
        this.$root.$emit("showSearch", this.q, false)
        this.$root.$emit("finishSearch")
      }
    }
  },
  mounted() {
    this.onQueryUpdate()
    this.$root.$on("search", (q: string) => {
      this.q = q
      this.refresh()
    })
  },
  computed: {
    pageCount() {
      // @ts-ignore
      if (this.data && this.data.total) {
        // @ts-ignore
        return Math.ceil(this.data.total / 10)
      }
      return 0
    }
  },
  watch: {
    "$route.query.q"() {
      this.q = <string>this.$route.query.q
      this.onQueryUpdate()
    },
    p() {
      this.search()
    },
    "$route.query.p"() {
      this.p = (this.$route.query.p && Number(this.$route.query.p) > 0) ? Number(this.$route.query.p) : 1
      this.onQueryUpdate()
    },
    "$route.query.i"() {
      this.i = (this.$route.query.i) ? [this.$route.query.i] : []
      this.onQueryUpdate()
    }
  }
});
</script>

<style>
.search-result {
  font-family: "Microsoft YaHei";
}

.search-result em {
  font-style: normal;
  color: #c00;
}

.link {
  color: #006d21;
  word-break: break-all;
  font-family: "Microsoft YaHei";
}

.content {
  font-family: "Microsoft YaHei";
  color: grey;
  word-break: break-all;
}

.search-result a:link {
  text-decoration: none;
  color: #001ba0;
  word-break: break-all;
}

.search-result a:hover {
  text-decoration: underline;
  color: #001ba0;
  word-break: break-all;
}
</style>
