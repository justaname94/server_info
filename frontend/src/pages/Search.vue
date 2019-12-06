<template>
  <div>
    <h2 class="text-center">Search Website's information</h2>
    <b-form @submit.prevent="getSiteInfo">
      <b-form-input
        v-model="url"
        :disabled="loading"
        placeholder="Ex: google.com"
        required
      ></b-form-input>
      <br />
      <div class="text-center">
        <b-button size="lg" type="submit" variant="primary">Search</b-button>
        <b-button size="lg" variant="primary" @click="navigateToHistory">
          History
        </b-button>
      </div>
    </b-form>
    <div class="text-center mt-3">
      <div v-if="loading">
        <b-spinner label="Loading..." class="mt-3"></b-spinner>
        <p>Please wait while we get the data, this could take a while...</p>
      </div>
    </div>
    <app-site-card
      v-if="Object.keys(website).length > 0"
      :website="website"
    ></app-site-card>
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
      url: "",
      website: {}
    };
  },
  methods: {
    getSiteInfo(evt) {
      this.website = {};
      this.loading = true;
      this.$http
        .get(`sites/${this.url}`)
        .then(response => {
          return response.json();
        })
        .then(data => {
          if (data.message !== undefined && data.message === "IN_PROGRESS") {
            // Recursively retry the request until is ready
            setTimeout(() => {
              this.getSiteInfo(evt);
            }, 3000);
          } else {
            this.loading = false;
            this.website = data;
          }
        });
    },
    navigateToHistory() {
      this.$router.push({ name: "history" });
    }
  }
};
</script>
