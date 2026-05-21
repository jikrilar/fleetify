import { api } from "../api.js";
import { createEl, formatDate, formatRupiah, showEmpty, statusBadge, statusLabel } from "../dom.js";
import { state } from "../state.js";

export async function renderHistoryPage(onDetail) {
  const tableBody = document.querySelector("#historyTableBody");
  state.reports = await api.getReports();

  if (state.reports.length === 0) {
    const row = createEl("tr");
    const cell = createEl("td", { colSpan: 8 });
    showEmpty(cell, "Riwayat laporan masih kosong.");
    row.appendChild(cell);
    tableBody.replaceChildren(row);
    return;
  }

  const fragment = document.createDocumentFragment();
  state.reports.forEach((report) => {
    const row = createEl("tr");
    row.appendChild(createEl("td", { text: `#${report.id}` }));
    row.appendChild(createEl("td", { text: report.user.username }));
    row.appendChild(createEl("td", { text: `${report.vehicle.license_plate} - ${report.vehicle.model}` }));
    row.appendChild(createEl("td", { text: `${report.odometer} km` }));
    const statusCell = createEl("td");
    statusCell.appendChild(statusBadge(report.status));
    row.appendChild(statusCell);
    row.appendChild(createEl("td", { text: formatRupiah(report.total_estimate) }));
    row.appendChild(createEl("td", { text: formatDate(report.created_at) }));
    row.appendChild(createEl("td", {}, [
      createEl("button", {
        className: "btn btn-outline-primary btn-sm",
        text: "Detail",
        type: "button",
        onClick: () => onDetail(report),
      }),
    ]));
    fragment.appendChild(row);
  });
  tableBody.replaceChildren(fragment);
}

export function exportReportsCsv() {
  const headers = ["ID", "SA", "Plat Nomor", "Model", "Odometer", "Keluhan", "Status", "Foto Awal", "Foto Bukti", "Total Estimasi", "Tanggal"];
  const rows = state.reports.map((report) => [
    report.id,
    report.user.username,
    report.vehicle.license_plate,
    report.vehicle.model,
    report.odometer,
    report.complaint,
    statusLabel(report.status),
    report.initial_photo || "",
    report.proof_photo || "",
    report.total_estimate,
    formatDate(report.created_at),
  ]);

  const csv = [headers, ...rows].map((row) => row.map(csvCell).join(",")).join("\n");
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = createEl("a", {
    href: url,
    download: "riwayat-laporan-fleetify.csv",
  });
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function csvCell(value) {
  const text = String(value ?? "");
  return `"${text.replaceAll('"', '""')}"`;
}
