<template>
  <b-container id="wrapper">
    <b-row>
      <b-col>
        <app-search></app-search>
      </b-col>
    </b-row>
    <!-- <app-site-card></app-site-card> -->
  </b-container>
</template>

<script>
import SiteCard from "./components/SiteCard";
import Search from "./pages/Search";

export default {
  name: "app",
  components: {
    AppSiteCard: SiteCard,
    AppSearch: Search
  },
  data() {
    return {
      msg: "Welcome to Your Vue.js App"
    };
  }
};
</script>

<style>
#wrapper {
  margin-top: 50px;
}
</style>