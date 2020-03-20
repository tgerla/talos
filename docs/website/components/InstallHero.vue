<template>
  <div class="w-3/4 mx-auto p-2 bg-gray-100 rounded-lg">
    <div class="flex">
      <div class="start-step">
        Install talosctl
      </div>

      <div class="flex-1 m-1 px-1 border-r-2 border-gray-400">
        <h4>On macOS:</h4>
        <pre><code class="language-bash">brew install talosctl</code></pre>
      </div>
      <div class="flex-1 m-1 px-1">
        <h4>On Linux:</h4>
        <p>Download the talosctl binary from GitHub.</p>
      </div>
    </div>

    <div class="flex">
      <div class="start-step">Start a cluster</div>
      <pre
        class="flex-1"
      ><code class="language-bash">talosctl cluster create</code></pre>
    </div>

    <div class="flex">
      <div class="start-step">Get your kubeconfig</div>
      <pre class="flex-1"><code class="language-bash">talosctl kubeconfig .
kubectl --kubeconfig kubeconfig config set-cluster talos_default --server https://127.0.0.1:6443</code></pre>
    </div>
    <div class="flex justify-center mx-auto">
      <a href="https://www.talos.dev/docs/v0.3/#v0.3/en/guides/getting-started">
        <button class="teal-cta-button">
          Read the Docs
        </button></a
      >
    </div>
  </div>
</template>

<style>
.start-step {
  @apply flex-none w-1/4 font-headings text-xl self-center;
}
</style>
