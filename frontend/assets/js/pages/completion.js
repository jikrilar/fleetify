import { api } from "../api.js";
import { createEl, showEmpty } from "../dom.js";
import { reportCard } from "../components.js";

export async function renderCompletePage(onChanged, onDetail) {
  const container = document.querySelector("#completeList");
  const reports = await api.getReports("APPROVED");

  if (reports.length === 0) {
    showEmpty(container, "Belum ada laporan yang sudah disetujui.");
    return;
  }

  const fragment = document.createDocumentFragment();
  reports.forEach((report) => {
    fragment.appendChild(completeCard(report, onChanged, onDetail));
  });
  container.replaceChildren(fragment);
}

function completeCard(report, onChanged, onDetail) {
  const card = reportCard(report, [{ label: "Detail", onClick: onDetail }]);
  const body = card.querySelector(".card-body");
  const form = createEl("form", { className: "mt-3 vstack gap-2" });
  const input = createEl("input", {
    className: "form-control form-control-sm",
    type: "text",
    placeholder: "contoh: bukti-selesai.jpg",
    required: true,
  });
  const button = createEl("button", {
    className: "btn btn-success btn-sm align-self-start",
    type: "submit",
    text: "Selesaikan",
  });

  form.appendChild(input);
  form.appendChild(button);
  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    button.disabled = true;
    button.textContent = "Memproses...";
    try {
      await api.completeReport(report.id, { proof_photo: input.value.trim() });
      await onChanged("Laporan berhasil diselesaikan");
    } finally {
      button.disabled = false;
      button.textContent = "Selesaikan";
    }
  });

  body.appendChild(form);
  return card;
}
