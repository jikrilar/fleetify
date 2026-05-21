import { createEl, formatDate, formatRupiah, statusBadge } from "./dom.js";

export function reportCard(report, actions = []) {
  const card = createEl("div", { className: "col-12 col-lg-6 col-xxl-4" });
  const body = createEl("div", { className: "card h-100" });
  const inner = createEl("div", { className: "card-body" });

  inner.appendChild(createEl("div", { className: "d-flex justify-content-between align-items-start gap-2 mb-2" }, [
    createEl("h2", { className: "card-title h6 mb-0", text: `Laporan #${report.id}` }),
    statusBadge(report.status),
  ]));
  inner.appendChild(createEl("p", { className: "mb-1", text: `${report.vehicle.license_plate} - ${report.vehicle.model}` }));
  inner.appendChild(createEl("p", { className: "text-secondary small mb-2", text: `Dibuat oleh ${report.user.username} pada ${formatDate(report.created_at)}` }));
  inner.appendChild(createEl("p", { className: "mb-2", text: report.complaint }));
  inner.appendChild(createEl("p", { className: "fw-semibold mb-3", text: `Total estimasi ${formatRupiah(report.total_estimate)}` }));

  const actionWrap = createEl("div", { className: "d-flex flex-wrap gap-2" });
  actions.forEach((action) => {
    actionWrap.appendChild(createEl("button", {
      className: action.className || "btn btn-outline-primary btn-sm",
      text: action.label,
      type: "button",
      onClick: () => action.onClick(report),
    }));
  });
  inner.appendChild(actionWrap);
  body.appendChild(inner);
  card.appendChild(body);
  return card;
}

export function detailContent(report) {
  const wrapper = createEl("div", { className: "vstack gap-3" });
  const summary = createEl("div", { className: "row g-3" });

  const fields = [
    ["ID", `#${report.id}`],
    ["SA", report.user.username],
    ["Kendaraan", `${report.vehicle.license_plate} - ${report.vehicle.model}`],
    ["Odometer", `${report.odometer} km`],
    ["Keluhan", report.complaint],
    ["Foto Awal", report.initial_photo || "-"],
    ["Foto Bukti", report.proof_photo || "-"],
    ["Total Estimasi", formatRupiah(report.total_estimate)],
    ["Tanggal Dibuat", formatDate(report.created_at)],
  ];

  fields.forEach(([label, value]) => {
    const col = createEl("div", { className: "col-md-6" });
    col.appendChild(createEl("div", { className: "text-secondary small", text: label }));
    col.appendChild(createEl("div", { className: "fw-semibold", text: value }));
    summary.appendChild(col);
  });

  const statusCol = createEl("div", { className: "col-md-6" });
  statusCol.appendChild(createEl("div", { className: "text-secondary small", text: "Status" }));
  statusCol.appendChild(statusBadge(report.status));
  summary.appendChild(statusCol);
  wrapper.appendChild(summary);

  const table = createEl("table", { className: "table table-sm" });
  const thead = createEl("thead");
  const headRow = createEl("tr");
  ["Item", "Tipe", "Qty", "Harga Snapshot", "Subtotal"].forEach((text) => {
    headRow.appendChild(createEl("th", { text }));
  });
  thead.appendChild(headRow);
  table.appendChild(thead);

  const tbody = createEl("tbody");
  report.items.forEach((item) => {
    const row = createEl("tr");
    row.appendChild(createEl("td", { text: item.master_item.item_name }));
    row.appendChild(createEl("td", { text: item.master_item.type }));
    row.appendChild(createEl("td", { text: String(item.quantity) }));
    row.appendChild(createEl("td", { text: formatRupiah(item.price_snapshot) }));
    row.appendChild(createEl("td", { text: formatRupiah(item.price_snapshot * item.quantity) }));
    tbody.appendChild(row);
  });
  table.appendChild(tbody);
  wrapper.appendChild(table);

  return wrapper;
}
