<template>
  <div>
    <div class="text-center">
      <h2 class="mb-3">History - View your past visited sites</h2>
      <b-button size="lg" variant="primary" @click="getHistory"
        >History</b-button
      >
      <b-button size="lg" variant="primary" @click="navigateToHome">
        Search
      </b-button>
      <div>
        <b-spinner v-if="loading" label="Loading..." class="mt-3"></b-spinner>
      </div>
    </div>
    <div v-for="website in websites" v-bind:key="website.domain">
      <app-site-card :website="website"></app-site-card>
    </div>
  </div>
</template>

<script>
import SiteCard from "../components/SiteCard";

export default {
  components: {
    AppSiteCard: SiteCard
  },
  data() {
    return {
      loading: false,
      websites: []
    };
  },

  methods: {
    getHistory(evt) {
      this.loading = true;
      this.$http
        .get("history")
        .then(response => {
          return response.json();
        })
        .then(data => {
          this.loading = false;
          this.websites = data;
        });
    },
    navigateToHome() {
      this.$router.push({ name: "home" });
    }
  }
};
</script>
