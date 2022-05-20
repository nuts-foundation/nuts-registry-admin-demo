<template>
  <div>
    <h1>Issue Verifiable Credential</h1>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <div class="space-y-4 w-full">
        <div v-if="feedbackMsg"
             :class="{ 'bg-green-300': responseState === 'success', 'bg-red-300': responseState === 'error'}"
             class="py-2 px-4 border rounded-md text-white">
          <svg v-if="responseState === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline"
               fill="none"
               viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
          <svg v-if="responseState === 'error'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline" fill="none"
               viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
          {{ feedbackMsg }}
        </div>

        <!-- VC to issue -->
        <div v-if="!issuedVC">
          <div>
            <label for="template-select">Choose a template</label>
            <select id="template-select" v-on:change="chooseTemplate">
              <option value="">Choose a template</option>
              <option :value="template.type" v-for="template in templates" :key="template.type">{{
                  template.type
                }}
              </option>
            </select>
          </div>
          <div>
            <label for="issuerdid-select">Issuer</label>
            <select id="issuerdid-select" v-model="vcToIssue.issuerDID">
              <option :value="issuer.did" v-for="issuer in availableIssuers" :key="issuer.did">{{ issuer.name }}:
                {{ issuer.did }}
              </option>
            </select>
          </div>
          <div>
            <label for="search-subject-input">Search for VC subject</label>
            <input type="text" v-model="subjectSearchQuery" autocomplete="false" id="search-subject-input"
                   v-on:input="searchSubjects" v-on:focusout="searchSubjects">
          </div>
          <div v-if="subjectSearchResults.length > 0" class="subject-search-results">
            <div v-for="(result, idx) in subjectSearchResults" :key="`result-${idx}`"
                 @click="selectSubject(result)">{{ result.organization.name }}
            </div>
          </div>
          <div>
            <label for="subjectDID-input">Issue VC to DID</label>
            <input id="subjectDID-input" v-model="vcToIssue.credentialSubjectDID" type="text">
          </div>

          <div>
            <label for="subject-textarea">Subject</label>
            <textarea id="subject-textarea" v-model="vcToIssue.credentialSubject"></textarea>
          </div>

          <div>
            <label for="publish-checkbox">Publish</label>
            <input id="publish-checkbox" v-model="vcToIssue.publishToNetwork" type="checkbox">
          </div>
          <div>
            <label for="visibility-select">Visibility</label>
            <select id="visibility-select" v-model="vcToIssue.visibility">
              <option>private</option>
              <option>public</option>
            </select>
          </div>
          <div>
            <label for="vccontext-input">Credential Context</label>
            <input id="vccontext-input" v-model="vcToIssue.vcContext" type="text">
          </div>
          <div>
            <label for="vctype-input">Credential Type</label>
            <input id="vctype-input" v-model="vcToIssue.vcType" type="text">
          </div>

          <div class="mt-4">
            <button id="issue-button" class="btn btn-primary" v-on:click="issueVC">Issue</button>
          </div>
        </div>
        <div v-if="issuedVC">
          <pre>{{JSON.stringify(issuedVC, null, 2)}}</pre>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  data() {
    return {
      fetchError: '',
      responseState: '',
      feedbackMsg: '',
      templates: [],
      availableIssuers: [],
      subjectSearchResults: [],
      issuedVC: null,
      subjectSearchQuery: null,
      vcToIssue: {
        credentialSubjectDID: null,
        credentialSubject: null,
        issuerDID: null,
        vcType: null,
        vcContext: null,
        publishToNetwork: true,
        visibility: "private",
      },
    }
  },
  mounted() {
    this.fetchTemplates()
    this.fetchIssuerDIDs()
  },
  methods: {
    selectSubject(vc) {
      this.vcToIssue.credentialSubjectDID = vc.subject
      this.subjectSearchQuery = ''
      this.subjectSearchResults = []
    },
    searchSubjects() {
      if (this.subjectSearchQuery === '') {
        this.subjectSearchResults = []
        return
      }
      this.$api.post('web/private/organizations', {name: this.subjectSearchQuery})
          .then(data => {
            this.subjectSearchResults = data
          })
          .catch(response => {
            this.fetchError = response.statusText
            this.subjectSearchResults = []
          })
    },
    fetchIssuerDIDs() {
      this.availableIssuers = []
      this.$api.get('web/private/service-provider')
          .then(responseData => {
            this.responseState = 'success'
            this.availableIssuers.push({did: responseData.id, name: "Service Provider"})
          })
          .catch(reason => {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
      this.$api.get('web/private/customers')
          .then(customers => {
            customers.forEach((customer) => {
              if (customer.active !== true) {
                return
              }
              this.availableIssuers.push({did: customer.did, name: customer.name})
            })
          })
          .catch(reason => {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    },
    fetchTemplates() {
      this.feedbackMsg = ''

      this.$api.get('web/private/vc/templates')
          .then(responseData => {
            this.templates = responseData
          })
          .catch(reason => {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    },
    chooseTemplate(el) {
      const templateType = el.target.value
      if (templateType === "") {
        return
      }
      const template = this.templates.filter((curr) => curr.type === templateType)[0]
      this.vcToIssue.vcType = template.type
      this.vcToIssue.vcContext = template.context
      this.vcToIssue.visibility = template.visibility
      this.vcToIssue.publishToNetwork = template.publishToNetwork
      this.vcToIssue.credentialSubject = JSON.stringify(template.credentialSubject, null, 2)
    },
    issueVC() {
      let inputCredentialSubject = JSON.parse(this.vcToIssue.credentialSubject)
      let credentialSubject = Object.assign({}, inputCredentialSubject) // copy
      credentialSubject.id = this.vcToIssue.credentialSubjectDID
      let request = {
        "@context": this.vcToIssue.vcContext,
        "issuer": this.vcToIssue.issuerDID,
        "type": this.vcToIssue.vcType,
        "credentialSubject": credentialSubject,
        "publishToNetwork": this.vcToIssue.publishToNetwork,
        "visibility": this.vcToIssue.visibility,
      };

      this.$api.post('web/private/vc', request)
          .then(responseData => {
            this.responseState = 'success'
            this.$emit('statusUpdate', 'Verifiable Credential Issued')
            this.feedbackMsg = ''
            this.issuedVC = responseData;
          })
          .catch(reason => {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    }
  }
}
</script>

<style>
textarea {
  height: 300px;
}

.subject-search-results {
  border-radius: 6px;
  border: 1px solid lightgray;
  padding: 5px;
  margin-top: 0;
}

.subject-search-results div {
  cursor: pointer;
}

.subject-search-results div:hover {
  background-color: #f1f1f1;
}
</style>