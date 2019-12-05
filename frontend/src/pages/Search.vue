<template>
  <div>
    <h2 class="text-center">Search Website's information</h2>
    <b-form @submit.prevent="getSiteInfo">
      <b-form-input v-model="url" placeholder="Ex: google.com"></b-form-input>
      <br />
      <div class="text-center">
        <b-button size="lg" type="submit" variant="primary">Search</b-button>
      </div>
    </b-form>
    <app-site-card v-if="Object.keys(website).length > 0" :website="website"></app-site-card>
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
      url: "",
      website: {}
    };
  },
  methods: {
    getSiteInfo(evt) {
      this.$http
        .get(`sites/${this.url}`)
        .then(response => {
          return response.json();
        })
        .then(data => {
          this.website = data;
        });
    }
  }
};
</script>
