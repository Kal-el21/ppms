import { useState } from "react";
import { PageHeader } from "../../../components/shared/PageHeader";
import ProfileSettings from "../components/ProfileSettings";
import NotificationSettings from "../components/NotificationSettings";
import SecuritySettings from "../components/SecuritySettings";

const tabs = [
  { key: "profile", label: "Profile" },
  { key: "notifications", label: "Notifikasi" },
  { key: "security", label: "Keamanan" },
];

export default function SettingsPage() {
  const [active, setActive] = useState("profile");

  return (
    <div>
      <PageHeader title="Settings" subtitle="Kelola profile, notifikasi, dan keamanan akun Anda" />

      <div className="flex gap-1 border-b border-border mb-6">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActive(tab.key)}
            className={`relative px-4 py-2.5 text-[13.5px] font-medium transition-colors ${
              active === tab.key
                ? "text-primary-700 dark:text-primary-400"
                : "text-ink-secondary hover:text-ink-primary"
            }`}
          >
            {tab.label}
            {active === tab.key && (
              <span className="absolute left-0 right-0 -bottom-px h-[2px] bg-primary-600 rounded-full" />
            )}
          </button>
        ))}
      </div>

      {active === "profile" && <ProfileSettings />}
      {active === "notifications" && <NotificationSettings />}
      {active === "security" && <SecuritySettings />}
    </div>
  );
}