# /*
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import devices by their ids
terraform import ome_devices.device_list_ids_default "<id1>,<id2>,<id3>"

# another way to import devices by their ids
terraform import ome_devices.device_list_ids "id:<id1>,<id2>,<id3>"

# import devices by their service tags
terraform import ome_devices.device_list_svc_tags "svc_tag:<svc_tag_1>,<svc_tag_2>,<svc_tag_3>"

# import all available devices
terraform import ome_devices.device_list_all ""
