import { api } from "../api.js";
import { reportCard } from "../components.js";
import { showEmpty } from "../dom.js";

export async function renderApprovalPage(onChanged, onDetail) {
  const container = document.querySelector("#approvalList");
  const reports = await api.getReports("PENDING_APPROVAL");

  if (reports.length === 0) {
    showEmpty(container, "Belum ada laporan yang menunggu approval.");
    return;
  }

  const fragment = document.createDocumentFragment();
  reports.forEach((report) => {
    fragment.appendChild(reportCard(report, [
      { label: "Detail", onClick: onDetail },
      {
        label: "Setujui",
        className: "btn btn-primary btn-sm",
        onClick: async (selectedReport) => {
          await api.approveReport(selectedReport.id);
          await onChanged("Laporan berhasil disetujui");
        },
      },
    ]));
  });
  container.replaceChildren(fragment);
}
