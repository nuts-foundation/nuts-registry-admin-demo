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

        <!-- Form -->
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
          <label for="issuerdid-input">Issuer</label>
          <input id="issuerdid-input" v-model="vcToIssue.issuerDID" type="text">
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
  },
  methods: {
    fetchTemplates() {
      this.feedbackMsg = ''

      this.$api.get('web/private/vc/templates')
          .then(responseData => {
            console.log(responseData)
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
      this.vcToIssue.credentialSubject = JSON.stringify(template.credentialSubject, null, 2)
    },
    issueVC() {
      /*


      {
    "@context":"https://kik-v.nl/context/v1.json",
    "issuer": "did:nuts:<authority DID>",
    "type": "ValidatedQueryCredential",
    "credentialSubject": {
        "id": "did:nuts:<data consumer DID>",
        "validatedQuery": {
            "profile": "https://kik-v2.gitlab.io/uitwisselprofielen/uitwisselprofiel-odb/",
            "ontology": "http://ontology.ontotext.com/publishing",
            "sparql": "PREFIX%20pub%3A%20%3Chttp%3A%2F%2Fontology.ontotext.com%2Ftaxonomy%2F%3E%0APREFIX%20publishing%3A%20%3Chttp%3A%2F%2Fontology.ontotext.com%2Fpublishing%23%3E%0ASELECT%20DISTINCT%20%3Fp%20%3FobjectLabel%20WHERE%20%7B%0A%20%20%20%20%3Chttp%3A%2F%2Fontology.ontotext.com%2Fresource%2Ftsk78dfdet4w%3E%20%3Fp%20%3Fo%20.%0A%20%20%20%20%7B%0A%20%20%20%20%20%20%20%20%3Fo%20pub%3AhasValue%20%3Fvalue%20.%0A%20%20%20%20%20%20%20%20%3Fvalue%20pub%3ApreferredLabel%20%3FobjectLabel%20.%0A%20%20%20%20%7D%20UNION%20%7B%0A%20%20%20%20%20%20%20%20%3Fo%20pub%3AhasValue%20%3FobjectLabel%20.%0A%20%20%20%20%20%20%20%20filter%20(isLiteral(%3FobjectLabel))%20.%0A%20%20%20%20%20%7D%0A%7D"
        }
    },
    "publishToNetwork": true,
    "visibility": "public"
}
       */
      let inputCredentialSubject = JSON.parse(this.vcToIssue.credentialSubject)
      console.log(inputCredentialSubject)
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
            // TODO: Show response data
            this.feedbackMsg = ''
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
</style>