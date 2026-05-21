export function createEl(tag, options = {}, children = []) {
  const element = document.createElement(tag);

  if (options.className) {
    element.className = options.className;
  }
  if (options.text !== undefined) {
    element.textContent = options.text;
  }
  if (options.type) {
    element.type = options.type;
  }
  if (options.value !== undefined) {
    element.value = options.value;
  }
  if (options.placeholder) {
    element.placeholder = options.placeholder;
  }
  if (options.id) {
    element.id = options.id;
  }
  if (options.required) {
    element.required = true;
  }
  if (options.min !== undefined) {
    element.min = options.min;
  }
  if (options.href) {
    element.href = options.href;
  }
  if (options.download) {
    element.download = options.download;
  }
  if (options.colSpan) {
    element.colSpan = options.colSpan;
  }
  if (options.dataset) {
    Object.entries(options.dataset).forEach(([key, value]) => {
      element.dataset[key] = value;
    });
  }
  if (options.onClick) {
    element.addEventListener("click", options.onClick);
  }
  if (options.onInput) {
    element.addEventListener("input", options.onInput);
  }

  children.forEach((child) => {
    element.appendChild(typeof child === "string" ? document.createTextNode(child) : child);
  });

  return element;
}

export function fillSelect(select, items, labelBuilder, valueBuilder = (item) => item.id) {
  const fragment = document.createDocumentFragment();
  items.forEach((item) => {
    fragment.appendChild(createEl("option", {
      value: valueBuilder(item),
      text: labelBuilder(item),
    }));
  });
  select.replaceChildren(fragment);
}

export function formatRupiah(value) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value || 0);
}

export function formatDate(value) {
  if (!value) {
    return "-";
  }
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export function statusLabel(status) {
  const labels = {
    PENDING_APPROVAL: "Menunggu Approval",
    APPROVED: "Disetujui",
    COMPLETED: "Selesai",
  };
  return labels[status] || status;
}

export function statusBadge(status) {
  const classes = {
    PENDING_APPROVAL: "badge bg-warning text-dark",
    APPROVED: "badge bg-primary",
    COMPLETED: "badge bg-success",
  };
  return createEl("span", {
    className: classes[status] || "badge bg-secondary",
    text: statusLabel(status),
  });
}

export function showEmpty(target, message) {
  target.replaceChildren(createEl("div", { className: "empty-state", text: message }));
}

export function showAlert(message, type = "success") {
  const alert = document.querySelector("#appAlert");
  alert.className = `alert alert-${type}`;
  alert.textContent = message;
  window.setTimeout(() => {
    alert.className = "alert d-none";
    alert.textContent = "";
  }, 4000);
}
