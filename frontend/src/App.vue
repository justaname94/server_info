<template>
  <b-container>
    <app-site-card></app-site-card>
  </b-container>
</template>

<script>
import SiteCard from "./components/SiteCard";

export default {
  name: "app",
  components: {
    AppSiteCard: SiteCard
  },
  data() {
    return {
      msg: "Welcome to Your Vue.js App"
    };
  }
};
</script>


