document.addEventListener("DOMContentLoaded", () => {
  const activeContainer = document.getElementById("active-container");
  const completedContainer = document.getElementById("completed-container");
  const activeEmpty = document.getElementById("active-empty");
  const completedEmpty = document.getElementById("completed-empty");

  const cards = document.querySelectorAll(".card");

  cards.forEach((card) => {
    const status = card.dataset.status;
    if (status === "pending") {
      setupCardSwipe(card, {
        activeContainer,
        completedContainer,
        activeEmpty,
        completedEmpty,
      });
      enableRestoreOnHold(card, {
        activeContainer,
        completedContainer,
        activeEmpty,
        completedEmpty,
      });
    } else {
      lockCard(card);
      moveToCompleted(card, status, {
        completedContainer,
        completedEmpty,
        activeContainer,
        activeEmpty,
      });
      enableRestoreOnHold(card, {
        activeContainer,
        completedContainer,
        activeEmpty,
        completedEmpty,
      });
    }
  });

  initKidTabs({
    activeContainer,
    completedContainer,
    activeEmpty,
    completedEmpty,
  });

  applyKidFilter(currentKidFilter, {
    activeContainer,
    completedContainer,
    activeEmpty,
    completedEmpty,
  });

  updateEmptyStates(activeContainer, completedContainer, activeEmpty, completedEmpty);
});

let currentKidFilter = null;
const LONG_PRESS_MS = 650;

const STATUS_TEXT = {
  pending: "Pending ⏳",
  done: "Done ✅",
  skipped: "Skipped ⏭",
};

function setupCardSwipe(card, containers) {
  if (card.dataset.swipeBound === "true") return;
  card.dataset.swipeBound = "true";

  let startX = 0;
  let currentX = 0;
  let isDragging = false;
  const threshold = 80; // px

  function getX(evt) {
    if (evt.touches && evt.touches.length > 0) {
      return evt.touches[0].clientX;
    }
    return evt.clientX;
  }

  function handleStart(evt) {
    if (card.classList.contains("locked")) return;
    isDragging = true;
    startX = getX(evt);
    currentX = startX;
    card.classList.add("dragging");
  }

  function handleMove(evt) {
    if (!isDragging) return;
    currentX = getX(evt);
    const deltaX = currentX - startX;

    // Move and slightly rotate card
    card.style.transform = `translateX(${deltaX}px) rotate(${deltaX / 25}deg)`;
  }

  function handleEnd(evt) {
    if (!isDragging) return;
    isDragging = false;

    const deltaX = currentX - startX;
    card.classList.remove("dragging");

    if (Math.abs(deltaX) > threshold) {
      const status = deltaX > 0 ? "done" : "skipped";
      swipeOff(card, deltaX, status);
    } else {
      // Snap back
      card.style.transform = "";
    }
  }

  function swipeOff(card, deltaX, status) {
    // animate off-screen
    const direction = deltaX > 0 ? 1 : -1;
    card.classList.add("settled", status);
    card.style.transform = `translateX(${direction * 400}px) rotate(${direction * 12}deg)`;

    updateStatusOnServer(card, status, containers);
  }

  // Mouse
  card.addEventListener("mousedown", (evt) => {
    evt.preventDefault();
    handleStart(evt);
  });

  // Touch
  card.addEventListener("touchstart", (evt) => {
    handleStart(evt);
  });

  window.addEventListener("mousemove", handleMove);
  window.addEventListener("touchmove", handleMove, { passive: true });

  window.addEventListener("mouseup", handleEnd);
  window.addEventListener("touchend", handleEnd);
}

function updateStatusOnServer(card, status, containers) {
  const taskId = parseInt(card.dataset.taskId, 10);

  fetch("/api/update_status", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ task_id: taskId, status }),
  })
    .then((resp) => resp.json())
    .then((data) => {
      if (!data.ok) {
        console.error("Server error updating status", data);
        // If you want, you can snap card back on error
        return;
      }

      // Update local status label + classes and move to correct stack
      card.dataset.status = status;
      card.classList.remove("pending", "done", "skipped");
      card.classList.add(status);

      const label = card.querySelector(".status-label");
      if (label) {
        label.textContent = STATUS_TEXT[status] || STATUS_TEXT.pending;
      }

      if (containers) {
        if (status === "pending") {
          moveToActive(card, containers);
        } else {
          moveToCompleted(card, status, containers);
        }
        applyKidFilter(currentKidFilter, containers);
      }
    })
    .catch((err) => {
      console.error("Network error updating status", err);
    });
}

function moveToCompleted(card, status, containers) {
  const { completedContainer, completedEmpty, activeContainer, activeEmpty } = containers;

  // Reset transform/drag visuals and lock the card
  card.style.transform = "";
  card.classList.remove("dragging", "settled");
  lockCard(card);

  // Ensure status class is set correctly
  card.classList.remove("pending", "done", "skipped");
  card.classList.add(status);

  // Move card into completed stack
  if (completedEmpty) {
    completedEmpty.hidden = true;
  }
  completedContainer.appendChild(card);

  updateEmptyStates(activeContainer, completedContainer, activeEmpty, completedEmpty);
}

function moveToActive(card, containers) {
  const { completedContainer, completedEmpty, activeContainer, activeEmpty } = containers;

  card.style.transform = "";
  card.classList.remove("dragging", "settled");
  unlockCard(card);

  card.classList.remove("pending", "done", "skipped");
  card.classList.add("pending");

  // Move card back to active stack
  if (activeEmpty) {
    activeEmpty.hidden = true;
  }
  activeContainer.appendChild(card);

  // Re-enable swipe interactions
  setupCardSwipe(card, { activeContainer, completedContainer, activeEmpty, completedEmpty });

  updateEmptyStates(activeContainer, completedContainer, activeEmpty, completedEmpty);
}

function lockCard(card) {
  card.classList.add("locked");
}

function unlockCard(card) {
  card.classList.remove("locked");
}

function updateEmptyStates(activeContainer, completedContainer, activeEmpty, completedEmpty) {
  const activeHasCards = countVisibleCards(activeContainer) > 0;
  const completedHasCards = countVisibleCards(completedContainer) > 0;

  if (activeEmpty) {
    activeEmpty.hidden = activeHasCards;
  }
  if (completedEmpty) {
    completedEmpty.hidden = completedHasCards;
  }
}

function enableRestoreOnHold(card, containers) {
  let pressTimer = null;

  const startPress = (evt) => {
    // Only allow restoring non-pending cards
    if (card.dataset.status === "pending") return;

    evt.preventDefault();
    clearTimeout(pressTimer);
    pressTimer = setTimeout(() => {
      const kid = card.querySelector(".card-kid")?.textContent || "";
      const title = card.querySelector(".card-title")?.textContent || "";
      const confirmed = window.confirm(
        `Mark "${title}"${kid ? ` for ${kid}` : ""} back to pending?`
      );
      if (confirmed) {
        updateStatusOnServer(card, "pending", containers);
      }
    }, LONG_PRESS_MS);
  };

  const cancelPress = () => {
    clearTimeout(pressTimer);
  };

  card.addEventListener("mousedown", startPress);
  card.addEventListener("touchstart", startPress);
  card.addEventListener("mouseup", cancelPress);
  card.addEventListener("mouseleave", cancelPress);
  card.addEventListener("touchend", cancelPress);
  card.addEventListener("touchcancel", cancelPress);
}

function applyKidFilter(kid, containers) {
  const cards = document.querySelectorAll(".card");
  const normalizedKid = kid || null;

  cards.forEach((card) => {
    const cardKid = card.dataset.kid || null;
    const hide = normalizedKid && cardKid !== normalizedKid;
    card.classList.toggle("is-hidden", hide);
  });

  if (containers) {
    const { activeContainer, completedContainer, activeEmpty, completedEmpty } = containers;
    updateEmptyStates(activeContainer, completedContainer, activeEmpty, completedEmpty);
  }
}

function initKidTabs(containers) {
  const tabStrip = document.getElementById("kid-tabs");
  if (!tabStrip) return;

  let kids = [];
  try {
    kids = JSON.parse(tabStrip.dataset.kids || "[]");
  } catch (err) {
    console.error("Unable to parse kid list", err);
  }

  const tabValues = ["All", ...kids];

  tabValues.forEach((kid, idx) => {
    const btn = document.createElement("button");
    btn.className = "tab-button";
    btn.textContent = kid === "All" ? "All" : kid;
    btn.dataset.kid = kid === "All" ? "" : kid;
    if (idx === 0) {
      btn.classList.add("active");
      currentKidFilter = null;
    }
    btn.addEventListener("click", () => {
      const targetKid = btn.dataset.kid || null;
      currentKidFilter = targetKid;
      setActiveTab(tabStrip, btn);
      applyKidFilter(currentKidFilter, containers);
    });
    tabStrip.appendChild(btn);
  });
}

function setActiveTab(strip, activeButton) {
  const buttons = strip.querySelectorAll(".tab-button");
  buttons.forEach((btn) => {
    btn.classList.toggle("active", btn === activeButton);
  });
}

function countVisibleCards(container) {
  return Array.from(container.querySelectorAll(".card")).filter(
    (card) => !card.classList.contains("is-hidden")
  ).length;
}
