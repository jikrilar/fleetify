import { api } from "../api.js";
import { createEl, fillSelect } from "../dom.js";
import { state } from "../state.js";

let itemRows;

export function initCreatePage(onSaved) {
  itemRows = document.querySelector("#itemRows");
  document.querySelector("#addItemButton").addEventListener("click", () => addItemRow());
  document.querySelector("#createReportForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const submitButton = document.querySelector("#createSubmitButton");
    submitButton.disabled = true;
    submitButton.textContent = "Menyimpan...";

    try {
      const payload = readForm();
      await api.createReport(payload);
      event.target.reset();
      renderCreatePage();
      await onSaved("Laporan berhasil dibuat");
    } finally {
      submitButton.disabled = false;
      submitButton.textContent = "Simpan Laporan";
    }
  });
}

export function renderCreatePage() {
  fillSelect(
    document.querySelector("#vehicleSelect"),
    state.vehicles,
    (vehicle) => `${vehicle.license_plate} - ${vehicle.model}`,
  );

  itemRows.replaceChildren();
  addItemRow();
}

function addItemRow() {
  const row = createEl("div", { className: "item-row row g-2 align-items-end" });
  const itemCol = createEl("div", { className: "col-md-7" });
  const qtyCol = createEl("div", { className: "col-md-3" });
  const actionCol = createEl("div", { className: "col-md-2" });

  const itemSelect = createEl("select", { className: "form-select item-select", required: true });
  fillSelect(itemSelect, state.items, (item) => `${item.item_name} - ${item.type}`);

  const qtyInput = createEl("input", {
    className: "form-control quantity-input",
    type: "number",
    min: "1",
    value: "1",
    required: true,
  });

  const removeButton = createEl("button", {
    className: "btn btn-outline-danger w-100",
    type: "button",
    text: "Hapus",
    onClick: () => {
      row.remove();
      if (itemRows.children.length === 0) {
        addItemRow();
      }
    },
  });

  itemCol.appendChild(createEl("label", { className: "form-label", text: "Item" }));
  itemCol.appendChild(itemSelect);
  qtyCol.appendChild(createEl("label", { className: "form-label", text: "Qty" }));
  qtyCol.appendChild(qtyInput);
  actionCol.appendChild(removeButton);
  row.appendChild(itemCol);
  row.appendChild(qtyCol);
  row.appendChild(actionCol);
  itemRows.appendChild(row);
}

function readForm() {
  const rows = Array.from(itemRows.querySelectorAll(".item-row"));
  return {
    vehicle_id: Number(document.querySelector("#vehicleSelect").value),
    odometer: Number(document.querySelector("#odometerInput").value),
    complaint: document.querySelector("#complaintInput").value.trim(),
    initial_photo: document.querySelector("#initialPhotoInput").value.trim(),
    items: rows.map((row) => ({
      item_id: Number(row.querySelector(".item-select").value),
      quantity: Number(row.querySelector(".quantity-input").value),
    })),
  };
}
