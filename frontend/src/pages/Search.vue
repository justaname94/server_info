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
        <b-button size="lg" variant="secondary" @click="navigateToHistory"
          >History</b-button
        >
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

    <b-modal ref="info-modal" hide-footer title="Request took too long">
      <div class="d-block text-center">
        <p>
          We have not been able to complete your request, please check back
          again in a few minutes
        </p>
      </div>
      <b-button class="mt-2" variant="outline-info" block @click="toggleModal"
        >Ok</b-button
      >
    </b-modal>
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
    sleep(ms) {
      return new Promise(resolve => setTimeout(resolve, ms));
    },
    toggleModal() {
      // We pass the ID of the button that we want to return focus to
      // when the modal has hidden
      this.$refs["info-modal"].toggle("#toggle-btn");
    },
    loadInfo() {
      this.$http
        .get(`sites/${this.url}`)
        .then(response => {
          return response.json();
        })
        .then(data => {
          if (data.message !== undefined && data.message === "IN_PROGRESS") {
            return;
          }
          this.loading = false;
          this.website = data;
        });
    },
    async getSiteInfo() {
      const tries = 6;
      const second = 1000;
      this.website = {};
      this.loading = true;

      for (let i = 0; i < tries; i++) {
        await this.loadInfo();
        if (!this.loading) {
          break;
        }
        await this.sleep(second << i);
      }
      if (this.loading) {
        this.loading = false;
        this.toggleModal();
      }
    },
    navigateToHistory() {
      this.$router.push({ name: "history" });
    }
  }
};
</script>
