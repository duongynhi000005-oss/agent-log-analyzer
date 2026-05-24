(function () {
  function initTooltips(root) {
    if (typeof window.tippy !== "function") return;
    const scope = root && typeof root.querySelectorAll === "function" ? root : document;
    const triggers = scope.querySelectorAll(".help-tip[data-tippy-content]:not([data-tooltip-ready])");
    if (triggers.length === 0) return;
    for (const trigger of triggers) {
      trigger.setAttribute("data-tooltip-ready", "true");
    }
    window.tippy(triggers, {
      allowHTML: false,
      appendTo: () => document.body,
      arrow: true,
      delay: [120, 60],
      duration: [120, 80],
      hideOnClick: true,
      interactive: true,
      maxWidth: 420,
      offset: [0, 8],
      placement: "top",
      theme: "agent-analyzer",
      trigger: "mouseenter focus click",
    });
  }

  window.AgentAnalyzerTooltips = { init: initTooltips };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", function () {
      initTooltips(document);
    });
  } else {
    initTooltips(document);
  }
})();
