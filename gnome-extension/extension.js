/* extension.js
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

/* exported init */
const {Gio, GObject} = imports.gi;

const QuickSettings = imports.ui.quickSettings;

const FeatureToggle = GObject.registerClass(
class FeatureToggle extends QuickSettings.QuickToggle {
    _init() {
        super._init({
            label: 'Feature Name',
            iconName: 'selection-mode-symbolic',
            toggleMode: true,
        });

        /* // Binding the toggle to a GSettings key
        this._settings = new Gio.Settings({
            schema_id: 'org.gnome.shell.extensions.example',
        });

        this._settings.bind('feature-enabled',
            this, 'checked',
            Gio.SettingsBindFlags.DEFAULT); */
    }
});

class Extension {
    constructor() {
        this._indicator = null;
    }

    enable() {
        this._indicator = new FeatureToggle();
    }

    disable() {
        this._indicator.destroy();
        this._indicator = null;
    }
}

function init() {
    return new Extension();
}
