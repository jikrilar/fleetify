import { api } from "./api.js";
import { detailContent } from "./components.js";
import { fillSelect, showAlert } from "./dom.js";
import { currentUser, setCurrentUser, state } from "./state.js";
import { initCreatePage, renderCreatePage } from "./pages/create-report.js";
import { renderApprovalPage } from "./pages/approval.js";
import { renderCompletePage } from "./pages/completion.js";
import { exportReportsCsv, renderHistoryPage } from "./pages/history.js";

const pages = {
  create: document.querySelector("#page-create"),
  approval: document.querySelector("#page-approval"),
  complete: document.querySelector("#page-complete"),
  history: document.querySelector("#page-history"),
};

document.addEventListener("DOMContentLoaded", start);

async function start() {
  bindNavigation();
  initCreatePage(afterWorkflowChanged);
  document.querySelector("#refreshApprovalButton").addEventListener("click", () => safeRenderCurrentPage());
  document.querySelector("#refreshCompleteButton").addEventListener("click", () => safeRenderCurrentPage());
  document.querySelector("#refreshHistoryButton").addEventListener("click", () => safeRenderCurrentPage());
  document.querySelector("#exportCsvButton").addEventListener("click", exportReportsCsv);

  try {
    state.users = await api.getUsers();
    setupUserSelect();
    await loadMasterData();
    applyRoleNavigation();
    renderCreatePage();
    await safeRenderCurrentPage();
  } catch (error) {
    showAlert(error.message, "danger");
  }
}

function bindNavigation() {
  document.querySelectorAll("[data-page]").forEach((button) => {
    button.addEventListener("click", async () => {
      state.currentPage = button.dataset.page;
      setActivePage();
      await safeRenderCurrentPage();
    });
  });
}

function setupUserSelect() {
	const select = document.querySelector("#userSelect");
	const savedUser = state.users.find((user) => String(user.id) === String(state.currentUserId));
	if ((!state.currentUserId || !savedUser) && state.users.length > 0) {
		setCurrentUser(state.users[0].id);
	}

  fillSelect(select, state.users, (user) => `${user.username} (${user.role})`);
  select.value = state.currentUserId;
  select.addEventListener("change", async () => {
    setCurrentUser(select.value);
    applyRoleNavigation();
    await safeRenderCurrentPage();
  });
}

async function loadMasterData() {
  state.vehicles = await api.getVehicles();
  state.items = await api.getItems();
}

function applyRoleNavigation() {
  const user = currentUser();
  const allowedPages = user?.role === "APPROVAL" ? ["approval", "history"] : ["create", "complete", "history"];

  document.querySelectorAll("[data-page]").forEach((button) => {
    const allowed = allowedPages.includes(button.dataset.page);
    button.classList.toggle("d-none", !allowed);
  });

  if (!allowedPages.includes(state.currentPage)) {
    state.currentPage = allowedPages[0];
  }

  setActivePage();
}

function setActivePage() {
  Object.entries(pages).forEach(([name, section]) => {
    section.classList.toggle("d-none", name !== state.currentPage);
  });

  document.querySelectorAll("[data-page]").forEach((button) => {
    button.classList.toggle("active", button.dataset.page === state.currentPage);
  });
}

async function safeRenderCurrentPage() {
  try {
    if (state.currentPage === "create") {
      renderCreatePage();
    }
    if (state.currentPage === "approval") {
      await renderApprovalPage(afterWorkflowChanged, showDetail);
    }
    if (state.currentPage === "complete") {
      await renderCompletePage(afterWorkflowChanged, showDetail);
    }
    if (state.currentPage === "history") {
      await renderHistoryPage(showDetail);
    }
  } catch (error) {
    showAlert(error.message, "danger");
  }
}

async function afterWorkflowChanged(message) {
  showAlert(message);
  await safeRenderCurrentPage();
}

function showDetail(report) {
  document.querySelector("#reportModalTitle").textContent = `Detail Laporan #${report.id}`;
  document.querySelector("#reportModalBody").replaceChildren(detailContent(report));
  const modal = new window.bootstrap.Modal(document.querySelector("#reportModal"));
  modal.show();
}
