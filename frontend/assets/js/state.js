const STORAGE_KEY = "fleetify_user_id";

export const state = {
  users: [],
  vehicles: [],
  items: [],
  reports: [],
  currentUserId: localStorage.getItem(STORAGE_KEY) || "",
  currentPage: "create",
};

export function setCurrentUser(id) {
  state.currentUserId = String(id);
  localStorage.setItem(STORAGE_KEY, state.currentUserId);
}

export function currentUser() {
  return state.users.find((user) => String(user.id) === String(state.currentUserId));
}
