<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-white text-black">
      <q-toolbar>
        <q-btn
          :to="{name:'index'}"
          flat
          no-caps>
          <q-icon name="search"/>
          <q-toolbar-title>
            {{ $t("title") }}
          </q-toolbar-title>
        </q-btn>
        <q-form>
          <q-input v-model="q" debounce="500" v-if="showSearch" class="q-ml-sm" dense outlined>
            <template v-slot:append>
              <q-btn dense :loading="searching" @click="doSearch" type="submit" round flat icon="search"/>
            </template>
          </q-input>
        </q-form>
        <q-space/>
        <div></div>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <router-view/>
    </q-page-container>
  </q-layout>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  name: 'MainLayout',
  data() {
    return {
      showSearch: false,
      q: "",
      searching: false
    }
  },
  methods: {
    doSearch() {
      this.searching = true
      this.$root.$emit("search", this.q)
    }
  },
  created() {
    this.$root.$on("finishSearch", () => {
      this.searching = false
    })
    this.$root.$on("showSearch", (q: string, show: boolean) => {
      this.showSearch = show
      this.q = q
    })
  },
});
</script>
