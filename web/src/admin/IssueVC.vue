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
          <select id="template-select" v-model="template">
            <option>private</option>
            <option>public</option>
          </select>
        </div>
        <div>
          <label for="vctype-input">Credential Type</label>
          <input id="vctype-input" v-model="vcToIssue.vcType" type="text">
        </div>
        <div>
          <label for="issuer-select">Issuer</label>
          <select id="issuer-select" v-model="vcToIssue.issuerDID">
            <option>private</option>
            <option>public</option>
          </select>
        </div>
        <div>
          <label for="subjectDID-input">Issue VC to DID</label>
          <input id="subjectDID-input" v-model="vcToIssue.credentialSubjectDID" type="text">
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
          <label for="subject-textarea">Subject</label>
          <textarea id="subject-textarea" v-model="vcToIssue.credentialSubject">

          </textarea>
        </div>

        <div class="mt-4">
          <button id="issue-button" class="btn btn-primary">Issue</button>
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
      vcToIssue: {
        credentialSubjectDID: null,
        credentialSubject: null,
        vcType: null,
        vcContext: null,
        publishToNetwork: true,
        visibility: "private",
      },
    }
  }
}
</script>

<style>
textarea {
  height: 300px;
}
</style>