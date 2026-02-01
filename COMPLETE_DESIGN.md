# Complete winx Library Design - IOCTL++ & Windows Forensic Analysis Framework

## Executive Summary

This design provides a **complete, from-scratch implementation** in pure Go (no external libraries except stdlib and `golang.org/x/sys/windows`) for:

1. **IOCTL++ Functionality** - Replicate IOCTL++ capabilities without requiring pre-compiled kernel drivers (IOCTLDump.sys, DevNameEnumWdm.sys)
2. **Windows Forensic Analysis** - Comprehensive artifact parsing, registry analysis, event log parsing, and filesystem forensics
3. **Network & PE Utilities** - DNS queries, packet analysis, PE inspection, and symbol handling

**Key Design Principles:**
- ✅ Pure Go implementation (build everything ourselves)
- ✅ No external dependencies beyond Go stdlib + `golang.org/x/sys/windows`
- ✅ User-mode and kernel-mode techniques
- ✅ Modular architecture for easy extension
- ✅ Complete API for IOCTL discovery, capture, replay, and fuzzing
- ✅ Comprehensive Windows forensic artifact parsing
- ✅ NTFS filesystem analysis and metadata extraction
- ✅ Browser forensics (Chromium, Mozilla, Safari)
- ✅ Registry hive parsing and event log analysis

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     winx IOCTL Framework                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────────┐  │
│  │   User-Mode    │  │  Kernel-Mode   │  │    Analysis      │  │
│  │    Hooking     │  │   Monitoring   │  │   & Utilities    │  │
│  └────────────────┘  └────────────────┘  └──────────────────┘  │
│         │                    │                     │            │
│  ┌──────▼──────┐      ┌──────▼──────┐      ┌──────▼──────┐    │
│  │ IAT Hooking │      │ ETW Tracing │      │ IOCTL Decode│    │
│  │Inline Hooks │      │ WMI Queries │      │   Fuzzing   │    │
│  │ DLL Inject  │      │Registry Mon.│      │ Capture/Rep │    │
│  └─────────────┘      └─────────────┘      └─────────────┘    │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │           Existing winx Core (device, service)            │  │
│  │  • DeviceIoControl • Driver Loading • Device Enumeration │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Complete Directory Structure

```
winx/
├── device/                    # Existing - Device I/O operations
│   ├── ioctl.go              # ✅ Exists: DeviceIoControl, OpenDevice
│   ├── setupdi.go            # ✅ Exists: Device enumeration
│   ├── constants.go          # ✅ Exists: IOCTL codes, CTL_CODE
│   ├── types.go              # ✅ Exists: Data structures
│   ├── device_test.go        # ✅ Exists: Tests
│   │
│   ├── decoder.go            # 🆕 NEW: IOCTL code decoding
│   ├── known_ioctls.go       # 🆕 NEW: Database of known IOCTLs
│   ├── discovery.go          # 🆕 NEW: Device path discovery helpers
│   ├── capture.go            # 🆕 NEW: Capture/replay file format
│   ├── fuzzer.go             # 🆕 NEW: IOCTL fuzzing engine
│   │
│   ├── usb.go                # 🆕 NEW: USB Storage Parser (usp)
│   ├── usb_types.go          # 🆕 NEW: USB device structures
│   ├── usb_registry.go       # 🆕 NEW: USBSTOR registry parsing
│   ├── usb_setupapi.go       # 🆕 NEW: SetupAPI log parsing
│   └── usb_timeline.go       # 🆕 NEW: USB connection timeline
│
├── service/                   # Existing - Driver management
│   ├── driver.go             # ✅ Exists: LoadDriver, UnloadDriver
│   └── driver_query.go       # 🆕 NEW: Query loaded drivers, status
│
├── hook/                      # 🆕 NEW: User-mode hooking framework
│   ├── iat_hook.go           # IAT (Import Address Table) hooking
│   ├── inline_hook.go        # Inline function hooking (detours)
│   ├── trampoline.go         # Trampoline generation for hooks
│   ├── hook_manager.go       # Unified hook management API
│   └── asm_x64.go            # x64 assembly helpers (jump generation)
│
├── inject/                    # 🆕 NEW: Process injection framework
│   ├── dll_inject.go         # DLL injection via CreateRemoteThread
│   ├── reflective_inject.go  # Reflective DLL injection
│   ├── thread_hijack.go      # Thread hijacking injection
│   └── hookdll/              # Hooking DLL to inject into processes
│       ├── main.go           # DLL entry point (DllMain)
│       ├── hooks.go          # Hook installation in injected DLL
│       └── ipc.go            # IPC back to host process
│
├── etw/                       # 🆕 NEW: ETW (Event Tracing for Windows)
│   ├── session.go            # ETW trace session management
│   ├── providers.go          # Kernel provider definitions
│   ├── consumer.go           # Event consumption and callbacks
│   └── kernel_events.go      # Kernel-specific event parsing
│
├── wmi/                       # 🆕 NEW: WMI queries (pure Go)
│   ├── query.go              # WMI query engine (COM-based)
│   ├── driver_info.go        # Win32_SystemDriver queries
│   └── device_info.go        # Win32_PnPEntity queries
│
├── registry/                  # 🆕 NEW: Registry Analysis (yaru)
│   ├── hive.go               # Registry hive parser
│   ├── types.go              # Registry structures
│   ├── cell.go               # Cell parsing (NK, VK, SK, LF, LH, RI, LI)
│   ├── value.go              # Value data parsing
│   ├── dirty.go              # Transaction log parsing
│   ├── monitor.go            # Registry change notifications
│   ├── driver_keys.go        # Driver registry key parsing
│   └── device_keys.go        # Device registry enumeration
│
├── memory/                    # 🆕 NEW: Memory manipulation
│   ├── process_memory.go     # Read/Write process memory
│   ├── protection.go         # VirtualProtect wrappers
│   └── allocation.go         # VirtualAlloc/Free wrappers
│
├── pe/                        # 🆕 NEW: PE parsing (for IAT hooking)
│   ├── parser.go             # Parse PE headers
│   ├── imports.go            # Parse import tables
│   └── exports.go            # Parse export tables
│   # Note: Extended PE utilities (scanner, viewer) in internal/pe/
│
├── ipc/                       # 🆕 NEW: Inter-process communication
│   ├── named_pipe.go         # Named pipe server/client
│   ├── shared_memory.go      # Shared memory sections
│   └── mailslot.go           # Mailslots for broadcasts
│
├── asm/                       # 🆕 NEW: Assembly code generation
│   ├── x64_encoder.go        # x64 instruction encoding
│   ├── jump_gen.go           # JMP instruction generation
│   └── disasm.go             # Basic disassembler (for hook length)
│
├── capture/                   # �� NEW: IOCTL capture system
│   ├── session.go            # Capture session management
│   ├── file_format.go        # .conf and .data file I/O
│   ├── buffer_log.go         # Buffer logging and replay
│   └── hook_bridge.go        # Bridge hooks to capture system
│
├── examples/                  # Example programs
│   ├── load_driver/          # Load and query driver
│   ├── hook_process/         # Hook DeviceIoControl in process
│   ├── capture_ioctls/       # Capture IOCTLs to files
│   ├── replay_ioctls/        # Replay from .conf/.data
│   ├── fuzz_driver/          # Fuzz driver with IOCTL range
│   └── monitor_system/       # System-wide IOCTL monitoring
│
└── tools/                     # Command-line tools
    ├── winxctl/              # CLI tool (like IOCTL++)
    │   ├── main.go
    │   ├── cmd_load.go       # Load driver command
    │   ├── cmd_hook.go       # Hook process command
    │   ├── cmd_capture.go    # Capture command
    │   ├── cmd_replay.go     # Replay command
    │   ├── cmd_fuzz.go       # Fuzz command
    │   └── cmd_monitor.go    # Monitor command
    └── hookdll.dll           # Compiled hooking DLL
│
├── internal/                      # 🆕 Internal implementation packages
│   │
│   ├── pe/                        # 🆕 PE Analysis (pe_view/pescan)
│   │   ├── parser.go              # PE file parser
│   │   ├── types.go               # PE structures
│   │   ├── dos.go                 # DOS header
│   │   ├── nt.go                  # NT headers
│   │   ├── sections.go            # Section parsing
│   │   ├── imports.go             # Import table
│   │   ├── exports.go             # Export table
│   │   ├── resources.go           # Resource parsing
│   │   ├── relocations.go         # Relocation table
│   │   ├── debug.go               # Debug directory
│   │   ├── tls.go                 # TLS directory
│   │   ├── security.go            # Security/signatures
│   │   ├── scanner.go             # PE anomaly scanner
│   │   └── viewer.go              # PE viewer utilities
│   │
│   ├── artifacts/                 # 🆕 Windows Artifact Parsers
│   │   ├── prefetch/              # Windows Prefetch Parser (pf)
│   │   │   ├── parser.go          # Prefetch file parser
│   │   │   ├── types.go           # Prefetch structures (v17-v30)
│   │   │   ├── decompress.go      # MAM compression handling
│   │   │   └── analysis.go        # Execution timeline analysis
│   │   │
│   │   ├── lnk/                   # Windows LNK Parsing Utility (lp)
│   │   │   ├── parser.go          # Shell link parser
│   │   │   ├── types.go           # LNK structures
│   │   │   ├── extradata.go       # Extra data blocks parsing
│   │   │   └── resolve.go         # Target resolution
│   │   │
│   │   ├── jumplist/              # Windows Jump List Parser (jmp)
│   │   │   ├── parser.go          # Jump list parser
│   │   │   ├── types.go           # Jump list structures
│   │   │   ├── destlist.go        # DestList stream parsing
│   │   │   └── olecf.go           # OLE Compound File parsing
│   │   │
│   │   ├── shellbag/              # Windows ShellBag Parser (sbag)
│   │   │   ├── parser.go          # ShellBag parser
│   │   │   ├── types.go           # Shell item structures
│   │   │   ├── itemid.go          # ITEMIDLIST parsing
│   │   │   └── bags.go            # BagMRU/Bags parsing
│   │   │
│   │   ├── shimcache/             # Windows AppCompatibility Cache (wacu)
│   │   │   ├── parser.go          # AppCompat Cache parser
│   │   │   ├── types.go           # ShimCache structures
│   │   │   └── registry.go        # Registry extraction
│   │   │
│   │   ├── shimdb/                # Windows Shim Database Parser (shims)
│   │   │   ├── parser.go          # SDB file parser
│   │   │   ├── types.go           # SDB structures
│   │   │   ├── tags.go            # TAG definitions
│   │   │   └── index.go           # Index parsing
│   │   │
│   │   ├── activitiescache/       # Timeline ActivitiesCache Parser (tac)
│   │   │   ├── parser.go          # ActivitiesCache.db parser
│   │   │   ├── types.go           # Activity structures
│   │   │   └── timeline.go        # Activity timeline generation
│   │   │
│   │   ├── indexdat/              # Windows 'index.dat' Parser (id)
│   │   │   ├── parser.go          # index.dat parser (legacy IE)
│   │   │   ├── types.go           # HASH table structures
│   │   │   └── url.go             # URL record parsing
│   │   │
│   │   ├── recycle/               # Trash Inspection & Analysis (tia)
│   │   │   ├── parser.go          # Recycle Bin parser
│   │   │   ├── info2.go           # INFO2 file parsing (XP)
│   │   │   ├── idollar.go         # $I/$R file parsing (Vista+)
│   │   │   └── analysis.go        # Deleted file analysis
│   │   │
│   │   ├── wpn/                   # Windows Push Notification DB Parser (wpn)
│   │   │   ├── parser.go          # WPN database parser
│   │   │   ├── types.go           # Notification structures
│   │   │   └── sqlite.go          # SQLite parsing
│   │   │
│   │   └── backstage/             # MS Office Backstage Parser (bs)
│   │       ├── parser.go          # Office MRU parser
│   │       ├── types.go           # Office structures
│   │       └── registry.go        # Office registry locations
│   │
│   ├── browser/                   # 🆕 Browser Artifact Parsers
│   │   ├── chromium/              # Chromium SQLite Parser (csp)
│   │   │   ├── parser.go          # Chromium artifact parser
│   │   │   ├── history.go         # History database
│   │   │   ├── cookies.go         # Cookies database
│   │   │   ├── downloads.go       # Downloads database
│   │   │   ├── logins.go          # Login Data (encrypted)
│   │   │   ├── bookmarks.go       # Bookmarks JSON
│   │   │   └── cache/             # Chromium Cache Parser (ccp)
│   │   │       ├── parser.go      # Cache parsing
│   │   │       ├── index.go       # Cache index
│   │   │       └── blockfile.go   # Block file format
│   │   │
│   │   ├── mozilla/               # Mozilla SQLite Parser (msp)
│   │   │   ├── parser.go          # Mozilla artifact parser
│   │   │   ├── places.go          # places.sqlite
│   │   │   ├── cookies.go         # cookies.sqlite
│   │   │   ├── formhistory.go     # formhistory.sqlite
│   │   │   ├── logins.go          # logins.json
│   │   │   └── cache/             # Mozilla Cache Parser (mcp)
│   │   │       ├── parser.go      # Cache2 parsing
│   │   │       ├── index.go       # Cache index
│   │   │       └── entries.go     # Cache entries
│   │   │
│   │   └── safari/                # Safari Artifact Parser (sap)
│   │       ├── parser.go          # Safari artifact parser
│   │       ├── history.go         # History.db
│   │       ├── downloads.go       # Downloads.plist
│   │       └── bookmarks.go       # Bookmarks.plist
│   │
│   ├── filesystem/                # 🆕 Filesystem Analysis
│   │   ├── ntfs/                  # NTFS Core Parsing
│   │   │   ├── volume.go          # NTFS volume handling
│   │   │   ├── boot.go            # Boot sector parsing
│   │   │   ├── types.go           # NTFS structures
│   │   │   └── cluster.go         # Cluster operations
│   │   │
│   │   ├── mft/                   # $MFT Parser (ntfswalk)
│   │   │   ├── parser.go          # MFT parser
│   │   │   ├── types.go           # MFT record structures
│   │   │   ├── record.go          # FILE record parsing
│   │   │   ├── attribute.go       # Attribute parsing
│   │   │   ├── filename.go        # $FILE_NAME attribute
│   │   │   ├── stdinfo.go         # $STANDARD_INFORMATION
│   │   │   ├── data.go            # $DATA attribute
│   │   │   ├── attrlist.go        # $ATTRIBUTE_LIST
│   │   │   └── runlist.go         # Data run parsing
│   │   │
│   │   ├── usnjrnl/               # Windows Journal Parser (jp)
│   │   │   ├── parser.go          # $UsnJrnl parser
│   │   │   ├── types.go           # USN record structures
│   │   │   ├── record.go          # USN record parsing
│   │   │   └── reasons.go         # Reason code definitions
│   │   │
│   │   ├── logfile/               # $LogFile Analysis (mala)
│   │   │   ├── parser.go          # $LogFile parser
│   │   │   ├── types.go           # Log structures
│   │   │   ├── restart.go         # Restart area parsing
│   │   │   ├── record.go          # Log record parsing
│   │   │   └── redo.go            # Redo/Undo operations
│   │   │
│   │   ├── indx/                  # Windows INDX Slack Parser (wisp)
│   │   │   ├── parser.go          # INDX buffer parser
│   │   │   ├── types.go           # INDEX structures
│   │   │   ├── entry.go           # Index entry parsing
│   │   │   ├── slack.go           # Slack space analysis
│   │   │   └── carver.go          # Deleted entry recovery
│   │   │
│   │   ├── fat/                   # FAT32 & exFAT Analysis (fata)
│   │   │   ├── fat32.go           # FAT32 parser
│   │   │   ├── exfat.go           # exFAT parser
│   │   │   ├── types.go           # FAT structures
│   │   │   ├── directory.go       # Directory entry parsing
│   │   │   └── recovery.go        # Deleted file recovery
│   │   │
│   │   ├── ntfsdir/               # NTFS Directory Enumerator (ntfsdir)
│   │   │   ├── enumerator.go      # Directory enumeration
│   │   │   ├── walker.go          # Recursive walker
│   │   │   └── filter.go          # File filtering
│   │   │
│   │   ├── ntfscopy/              # NTFS File Copy Utility (ntfscopy)
│   │   │   ├── copy.go            # Raw NTFS file copy
│   │   │   ├── ads.go             # Alternate Data Streams
│   │   │   └── locked.go          # Locked file handling
│   │   │
│   │   └── gena/                  # Graphical Engine for NTFS Analysis
│   │       ├── engine.go          # Analysis engine
│   │       ├── visualize.go       # Data visualization
│   │       └── export.go          # Export utilities
│   │
│   ├── evtx/                      # 🆕 Windows Event Log Parser (evtwalk/evtx_view)
│   │   ├── parser.go              # EVTX file parser
│   │   ├── types.go               # EVTX structures
│   │   ├── chunk.go               # Chunk parsing
│   │   ├── record.go              # Event record parsing
│   │   ├── binxml.go              # Binary XML parsing
│   │   ├── template.go            # Template definitions
│   │   ├── message.go             # Message table resolution
│   │   └── fragment/              # EVTX Fragment Parser (evtfx)
│   │       ├── parser.go          # Fragment recovery
│   │       └── carver.go          # Fragment carving
│   │
│   ├── messagetable/              # 🆕 Event Log MessageTables Offline (elmo)
│   │   ├── parser.go              # Message table parser
│   │   ├── resource.go            # PE resource extraction
│   │   └── format.go              # FormatMessage recreation
│   │
│   ├── cafae/                     # 🆕 Computer Account Artifact Extractor
│   │   ├── extractor.go           # Artifact extraction
│   │   ├── sam.go                 # SAM hive parsing
│   │   ├── security.go            # SECURITY hive parsing
│   │   └── software.go            # SOFTWARE hive parsing
│   │
│   ├── tela/                      # 🆕 Trace Event Log and Analysis
│   │   ├── analyzer.go            # Event correlation
│   │   ├── timeline.go            # Event timeline
│   │   └── report.go              # Analysis reports
│   │
│   ├── network/                   # 🆕 Network Support Utilities
│   │   ├── dns/                   # DNS Query Utility (dqu)
│   │   │   ├── query.go           # DNS query implementation
│   │   │   ├── types.go           # DNS structures
│   │   │   ├── resolver.go        # Custom resolver
│   │   │   ├── cache.go           # DNS cache inspection
│   │   │   └── record.go          # Record type parsing
│   │   │
│   │   ├── pcap/                  # Packet Capture ICMP Carver (pic)
│   │   │   ├── parser.go          # PCAP/PCAPNG parser
│   │   │   ├── types.go           # Packet structures
│   │   │   ├── icmp.go            # ICMP packet parsing
│   │   │   ├── carver.go          # ICMP data carving
│   │   │   └── writer.go          # PCAP writing
│   │   │
│   │   ├── netxfer/               # Network Xfer Client/Server (nx)
│   │   │   ├── server.go          # Transfer server
│   │   │   ├── client.go          # Transfer client
│   │   │   ├── protocol.go        # Transfer protocol
│   │   │   └── crypto.go          # Optional encryption
│   │   │
│   │   └── minx/                  # Modular Inspection Network Xfer Agent
│   │       ├── agent.go           # MINX agent
│   │       ├── modules.go         # Inspection modules
│   │       ├── transport.go       # Transport layer
│   │       └── command.go         # Command handling
│   │
│   ├── vss/                       # 🆕 Volume Shadow Snapshot (vssenum)
│   │   ├── enumerator.go          # VSS enumeration
│   │   ├── types.go               # VSS structures
│   │   ├── snapshot.go            # Snapshot access
│   │   ├── diff.go                # Snapshot diffing
│   │   └── copy.go                # Copy from snapshots
│   │
│   ├── symbol/                    # 🆕 Windows Symbol Fetch Utility (sf)
│   │   ├── fetch.go               # Symbol downloader
│   │   ├── server.go              # Symbol server protocol
│   │   ├── pdb.go                 # PDB parsing basics
│   │   ├── cache.go               # Symbol cache management
│   │   └── guid.go                # GUID/age extraction
│   │
│   ├── csvdx/                     # 🆕 CSV Data eXchange
│   │   ├── parser.go              # CSV parser
│   │   ├── writer.go              # CSV writer
│   │   ├── transform.go           # Data transformation
│   │   └── merge.go               # CSV merging
│   │
│   └── disk/                      # 🆕 Disk Utility & Packer (dup)
│       ├── disk.go                # Raw disk access
│       ├── partition.go           # Partition parsing
│       ├── mbr.go                 # MBR parsing
│       ├── gpt.go                 # GPT parsing
│       ├── packer.go              # Disk image packing
│       └── imaging.go             # Disk imaging
│
└── examples/                      # Extended example programs
    ├── load_driver/               # Load and query driver
    ├── hook_process/              # Hook DeviceIoControl
    ├── capture_ioctls/            # Capture IOCTLs to files
    ├── replay_ioctls/             # Replay from .conf/.data
    ├── fuzz_driver/               # Fuzz driver with IOCTL range
    ├── monitor_system/            # System-wide monitoring
    ├── parse_prefetch/            # Parse prefetch files (pf)
    ├── parse_lnk/                 # Parse LNK files (lp)
    ├── parse_jumplist/            # Parse jump lists (jmp)
    ├── parse_usb/                 # Parse USB history (usp)
    ├── parse_shellbag/            # Parse ShellBags (sbag)
    ├── parse_shimcache/           # Parse AppCompat cache (wacu)
    ├── parse_shimdb/              # Parse SDB files (shims)
    ├── parse_activities/          # Parse ActivitiesCache (tac)
    ├── parse_indexdat/            # Parse index.dat (id)
    ├── parse_recycle/             # Parse Recycle Bin (tia)
    ├── parse_wpn/                 # Parse WPN database (wpn)
    ├── parse_backstage/           # Parse Office MRU (bs)
    ├── parse_chromium/            # Parse Chromium artifacts (csp)
    ├── parse_chromium_cache/      # Parse Chromium cache (ccp)
    ├── parse_mozilla/             # Parse Mozilla artifacts (msp)
    ├── parse_mozilla_cache/       # Parse Mozilla cache (mcp)
    ├── parse_safari/              # Parse Safari artifacts (sap)
    ├── parse_evtx/                # Parse EVTX files (evtx_view/evtwalk)
    ├── parse_evtx_fragment/       # Parse EVTX fragments (evtfx)
    ├── parse_messagetable/        # Parse message tables (elmo)
    ├── extract_cafae/             # Extract account artifacts (cafae)
    ├── analyze_tela/              # Trace event analysis (tela)
    ├── parse_registry/            # Parse registry hives (yaru)
    ├── parse_mft/                 # Parse $MFT (ntfswalk)
    ├── parse_usnjrnl/             # Parse $UsnJrnl (jp)
    ├── parse_logfile/             # Parse $LogFile (mala)
    ├── parse_indx/                # Parse INDX slack (wisp)
    ├── parse_fat/                 # Parse FAT filesystems (fata)
    ├── ntfs_enumerate/            # NTFS directory enumeration (ntfsdir)
    ├── ntfs_copy/                 # NTFS file copy (ntfscopy)
    ├── ntfs_analyze/              # NTFS graphical analysis (gena)
    ├── dns_query/                 # DNS queries (dqu)
    ├── pcap_carve/                # ICMP carving (pic)
    ├── netxfer/                   # Network transfer (nx)
    ├── minx_agent/                # MINX agent (minx)
    ├── pe_view/                   # PE viewer (pe_view)
    ├── pe_scan/                   # PE scanner (pescan)
    ├── vss_enumerate/             # VSS enumeration (vssenum)
    ├── symbol_fetch/              # Symbol fetching (sf)
    ├── csv_exchange/              # CSV utilities (csvdx)
    └── disk_util/                 # Disk utilities (dup)
```

---

## Complete Functionality Matrix

### 1. Core Device & IOCTL Operations (Existing + Extensions)

#### Existing Functions (✅ Already in winx)
```
device/ioctl.go:
  ✅ CreateFile(fileName, access, shareMode, ...) -> Handle
  ✅ CloseHandle(handle) -> bool
  ✅ DeviceIoControl(device, code, in, out, ...) -> bool
  ✅ DeviceIoControlBytes(device, code, inBuf, outSize) -> []byte
  ✅ ReadFile(handle, buffer, ...) -> bool
  ✅ WriteFile(handle, buffer, ...) -> bool
  ✅ OpenDevice(path, access) -> Handle
  ✅ OpenDeviceReadWrite(path) -> Handle
  ✅ QueryDosDevice(deviceName) -> []string
  ✅ FindSymbolicLinksByPattern(pattern) -> map[string][]string

device/setupdi.go:
  ✅ SetupDiGetClassDevs(guid, flags) -> Handle
  ✅ SetupDiEnumDeviceInterfaces(...) -> bool
  ✅ SetupDiGetDeviceInterfaceDetail(...) -> string
  ✅ EnumerateDevices(guid, flags) -> []string
  ✅ EnumerateDevicesWithInfo(flags) -> []DeviceInfo
  ✅ FindDevicesByService(serviceName) -> []DeviceInfo
  ✅ GetDriverDevicePaths(serviceName) -> []string

device/constants.go:
  ✅ CTL_CODE(deviceType, function, method, access) -> uint32
  ✅ IOCTL_DISK_GET_DRIVE_GEOMETRY (constant)
  ✅ IOCTL_STORAGE_QUERY_PROPERTY (constant)
  ✅ FILE_DEVICE_* constants
  ✅ METHOD_* constants

service/driver.go:
  ✅ OpenSCManager(machine, database, access) -> Handle
  ✅ CreateService(scm, name, displayName, ...) -> Handle
  ✅ OpenService(scm, name, access) -> Handle
  ✅ StartService(service, args) -> bool
  ✅ ControlService(service, control, status) -> bool
  ✅ DeleteService(service) -> bool
  ✅ QueryServiceStatus(service, status) -> bool
  ✅ LoadDriver(path, name) -> Handle
  ✅ UnloadDriver(service) -> error
```

#### New Functions (🆕 To Implement)

**device/decoder.go** - IOCTL Code Decoding
```go
🆕 DecodeIOCTL(ioctlCode uint32) -> IOCTLComponents
🆕 FormatIOCTL(ioctlCode uint32) -> string
🆕 GetDeviceTypeName(deviceType uint32) -> string
🆕 GetMethodName(method uint32) -> string
🆕 GetAccessName(access uint32) -> string
🆕 ParseIOCTLString(ioctlStr string) -> uint32, error
🆕 ValidateIOCTL(ioctlCode uint32) -> error

// IOCTLComponents structure
type IOCTLComponents struct {
    IOCTLCode      uint32
    DeviceType     uint32
    DeviceTypeName string
    Function       uint32
    Method         uint32
    MethodName     string
    Access         uint32
    AccessName     string
}
```

**device/known_ioctls.go** - Known IOCTL Database
```go
🆕 GetKnownIOCTLs() -> map[uint32]string
🆕 LookupIOCTL(ioctlCode uint32) -> (string, bool)
🆕 RegisterIOCTL(ioctlCode uint32, name string)
🆕 LoadIOCTLDatabase(filepath string) -> error
🆕 SaveIOCTLDatabase(filepath string) -> error
🆕 GetIOCTLsByDeviceType(deviceType uint32) -> map[uint32]string
🆕 SearchIOCTLByName(pattern string) -> []uint32
```

**device/discovery.go** - Enhanced Device Discovery
```go
🆕 DetectNewDevices(before, after []DeviceInfo) -> []DeviceInfo
🆕 FindDeviceDelta(baseline, current []DeviceInfo) -> DeviceDelta
🆕 GetAllDevicePaths() -> []string
🆕 GetDevicesByInterface(interfaceGuid GUID) -> []DeviceInfo
🆕 ResolveDevicePath(symbolicLink string) -> (string, error)
🆕 GetDeviceProperties(devicePath string) -> DeviceProperties

type DeviceDelta struct {
    Added   []DeviceInfo
    Removed []DeviceInfo
    Changed []DeviceInfo
}
```

**device/capture.go** - IOCTL Capture & Replay
```go
🆕 SaveCaptureConfig(config *CaptureConfig, filepath string) -> error
🆕 LoadCaptureConfig(filepath string) -> (*CaptureConfig, error)
🆕 SaveCaptureData(data *CaptureData, filepath string) -> error
🆕 LoadCaptureData(filepath string) -> (*CaptureData, error)
🆕 ReplayIOCTL(configPath, dataPath string) -> (*IOCTLResponse, error)
🆕 ReplayIOCTLModified(config, data, modifyFunc) -> (*IOCTLResponse, error)
🆕 ReplaySession(sessionDir string) -> ([]IOCTLResponse, error)

type CaptureConfig struct {
    DevicePath      string
    DeviceName      string
    IOCTLCode       uint32
    IOCTLDecoded    *IOCTLComponents
    InputSize       uint32
    OutputSize      uint32
    Timestamp       time.Time
    ProcessName     string
    ProcessID       uint32
}

type CaptureData struct {
    InputBuffer  []byte
    OutputBuffer []byte
}
```

**device/fuzzer.go** - IOCTL Fuzzing Engine
```go
🆕 FuzzIOCTL(devicePath string, opts *FuzzOptions) -> (*FuzzResults, error)
🆕 FuzzIOCTLRange(device, startCode, endCode) -> (*FuzzResults, error)
🆕 GenerateIOCTLRange(deviceType, startFunc, endFunc, method, access) -> []uint32
🆕 GenerateAllMethodVariants(deviceType, function, access) -> []uint32
🆕 DiscoverValidIOCTLs(devicePath string, deviceType uint32) -> ([]uint32, error)
🆕 TestIOCTLCode(device, code, input, outputSize) -> (*IOCTLTestResult, error)

type FuzzOptions struct {
    StartCode       uint32
    EndCode         uint32
    DeviceType      uint32
    Method          uint32
    InputData       []byte
    OutputSize      uint32
    Timeout         time.Duration
    OnSuccess       func(code uint32, response []byte)
    OnError         func(code uint32, err error)
    Parallel        int
}

type FuzzResults struct {
    TotalTried      uint32
    SuccessCount    uint32
    ErrorCount      uint32
    TimeoutCount    uint32
    SuccessfulCodes []uint32
    ErrorCodes      map[uint32]error
    Timing          map[uint32]time.Duration
}
```

**service/driver_query.go** - Driver Status Queries
```go
🆕 GetLoadedDrivers() -> ([]DriverInfo, error)
🆕 IsDriverRunning(serviceName string) -> (bool, error)
🆕 GetDriverStatus(serviceName string) -> (*DriverStatus, error)
🆕 GetDriverPath(serviceName string) -> (string, error)
🆕 EnumerateDriverServices() -> ([]string, error)

type DriverInfo struct {
    ServiceName  string
    DisplayName  string
    DriverPath   string
    Status       uint32
    StartType    uint32
    ErrorControl uint32
}

type DriverStatus struct {
    ServiceName    string
    CurrentState   uint32
    ControlsAccepted uint32
    Win32ExitCode  uint32
}
```

---

### 2. User-Mode Hooking Framework (🆕 New Package: hook/)

**hook/hook_manager.go** - Unified Hook Management
```go
🆕 NewHookManager() -> *HookManager
🆕 InstallHook(target, hookFunc, hookType) -> (HookHandle, error)
🆕 RemoveHook(handle HookHandle) -> error
🆕 RemoveAllHooks() -> error
🆕 GetInstalledHooks() -> []HookInfo
🆕 EnableHook(handle HookHandle) -> error
🆕 DisableHook(handle HookHandle) -> error

type HookManager struct {
    hooks map[HookHandle]*Hook
    mutex sync.RWMutex
}

type HookType int
const (
    HookTypeIAT HookType = iota
    HookTypeInline
    HookTypeVTable
)
```

**hook/iat_hook.go** - IAT Hooking
```go
🆕 NewIATHook(moduleName, functionName string) -> *IATHook
🆕 Install(hookFunction uintptr) -> error
🆕 Remove() -> error
🆕 GetOriginalFunction() -> uintptr
🆕 FindIATEntry(module, function) -> (uintptr, error)
🆕 PatchIAT(iatEntry, hookAddr uintptr) -> error

type IATHook struct {
    ModuleName      string
    FunctionName    string
    IATEntry        uintptr
    OriginalPointer uintptr
    HookPointer     uintptr
    IsInstalled     bool
}
```

**hook/inline_hook.go** - Inline Function Hooking
```go
🆕 NewInlineHook(targetAddr, hookAddr uintptr) -> *InlineHook
🆕 Install() -> error
🆕 Remove() -> error
🆕 GetTrampoline() -> uintptr
🆕 BuildTrampoline(originalBytes []byte) -> (uintptr, error)
🆕 CalculateHookLength(targetAddr uintptr) -> (int, error)

type InlineHook struct {
    TargetAddress   uintptr
    HookAddress     uintptr
    OriginalBytes   []byte
    TrampolineAddr  uintptr
    HookLength      int
    IsInstalled     bool
}
```

**hook/trampoline.go** - Trampoline Generation
```go
🆕 AllocateTrampoline(size uintptr) -> (uintptr, error)
🆕 BuildTrampolineCode(originalBytes []byte, returnAddr uintptr) -> []byte
🆕 FreeTrampoline(addr uintptr) -> error
🆕 DisassembleForTrampoline(addr uintptr, minLength int) -> ([]byte, error)
```

**hook/asm_x64.go** - x64 Assembly Helpers
```go
🆕 EncodeJmpRel32(target, source uintptr) -> []byte
🆕 EncodeJmpAbs64(target uintptr) -> []byte
🆕 EncodeCallRel32(target, source uintptr) -> []byte
🆕 EncodePush64(value uint64) -> []byte
🆕 EncodePop(register byte) -> []byte
🆕 EncodeMovRAX(value uint64) -> []byte
🆕 GetInstructionLength(addr uintptr) -> (int, error)
```

---

### 3. Process Injection Framework (🆕 New Package: inject/)

**inject/dll_inject.go** - Standard DLL Injection
```go
🆕 InjectDLL(processID uint32, dllPath string) -> error
🆕 InjectDLLEx(pid uint32, dllPath string, method InjectionMethod) -> error
🆕 EjectDLL(pid uint32, dllPath string) -> error
🆕 GetInjectedModules(pid uint32) -> ([]string, error)

type InjectionMethod int
const (
    MethodCreateRemoteThread InjectionMethod = iota
    MethodThreadHijack
    MethodReflective
    MethodQueueUserAPC
)
```

**inject/reflective_inject.go** - Reflective DLL Injection
```go
🆕 InjectReflective(pid uint32, dllBytes []byte) -> error
🆕 LoadRemoteDLL(pid uint32, dllData []byte) -> (uintptr, error)
🆕 ResolveImports(pid uint32, baseAddr uintptr, pe *PE) -> error
🆕 RelocateImage(pid uint32, baseAddr uintptr, pe *PE) -> error
```

**inject/thread_hijack.go** - Thread Hijacking
```go
🆕 InjectViaThreadHijack(pid uint32, dllPath string) -> error
🆕 SuspendThread(threadID uint32) -> error
🆕 GetThreadContext(thread Handle) -> (*CONTEXT, error)
🆕 SetThreadContext(thread Handle, ctx *CONTEXT) -> error
🆕 HijackThread(threadID uint32, shellcode []byte) -> error
```

**inject/hookdll/** - Hooking DLL (Compiled to .dll)
```go
// hookdll/main.go
🆕 DllMain(hinstDLL, fdwReason, lpvReserved) -> bool

// hookdll/hooks.go
🆕 InstallDeviceIoControlHook() -> error
🆕 DeviceIoControlDetour(params) -> result
🆕 LogIOCTL(code, device, in, out) -> error

// hookdll/ipc.go
🆕 ConnectToHost() -> error
🆕 SendCapturedIOCTL(capture *CaptureData) -> error
🆕 ReceiveCommands() -> (*Command, error)
```

---

### 4. ETW Framework (🆕 New Package: etw/)

**etw/session.go** - ETW Session Management
```go
🆕 NewSession(sessionName string) -> *Session
🆕 Start() -> error
🆕 Stop() -> error
🆕 EnableProvider(guid GUID, level, keywords uint64) -> error
🆕 DisableProvider(guid GUID) -> error
🆕 GetSessionInfo() -> (*SessionInfo, error)

type Session struct {
    Name          string
    Handle        Handle
    Properties    *EVENT_TRACE_PROPERTIES
    IsRunning     bool
}
```

**etw/providers.go** - Kernel Provider Definitions
```go
🆕 GetKernelFileProvider() -> GUID
🆕 GetKernelDiskProvider() -> GUID
🆕 GetKernelNetworkProvider() -> GUID
🆕 GetKernelProcessProvider() -> GUID
🆕 ListAvailableProviders() -> ([]ProviderInfo, error)

type ProviderInfo struct {
    GUID        GUID
    Name        string
    Description string
}
```

**etw/consumer.go** - Event Consumption
```go
🆕 NewConsumer(sessionName string) -> *Consumer
🆕 StartConsuming(callback EventCallback) -> error
🆕 StopConsuming() -> error
🆕 ProcessEvents() -> error

type EventCallback func(event *EVENT_RECORD)

type Consumer struct {
    SessionName string
    TraceHandle Handle
    Callback    EventCallback
}
```

**etw/kernel_events.go** - Kernel Event Parsing
```go
🆕 ParseFileIOEvent(event *EVENT_RECORD) -> (*FileIOEvent, error)
🆕 ParseDiskIOEvent(event *EVENT_RECORD) -> (*DiskIOEvent, error)
🆕 ParseProcessEvent(event *EVENT_RECORD) -> (*ProcessEvent, error)
🆕 ExtractIOCTLFromEvent(event *EVENT_RECORD) -> (uint32, bool)
```

---

### 5. WMI Query Framework (🆕 New Package: wmi/)

**wmi/query.go** - WMI Query Engine (Pure Go, COM-based)
```go
🆕 NewWMIClient() -> (*Client, error)
🆕 Query(wqlQuery string) -> (*ResultSet, error)
🆕 ExecQuery(query string, callback RowCallback) -> error
🆕 Close() -> error

type Client struct {
    locator  *IWbemLocator
    service  *IWbemServices
}

type ResultSet struct {
    Columns []string
    Rows    []map[string]interface{}
}
```

**wmi/driver_info.go** - Driver Queries
```go
🆕 GetSystemDrivers() -> ([]SystemDriver, error)
🆕 GetLoadedDrivers() -> ([]SystemDriver, error)
🆕 GetDriverByName(name string) -> (*SystemDriver, error)
🆕 GetDriverDependencies(name string) -> ([]string, error)

type SystemDriver struct {
    Name        string
    DisplayName string
    PathName    string
    State       string
    Started     bool
    ServiceType string
    StartMode   string
}
```

**wmi/device_info.go** - Device Queries
```go
🆕 GetPnPDevices() -> ([]PnPDevice, error)
🆕 GetDevicesByClass(classGuid string) -> ([]PnPDevice, error)
🆕 GetDevicesByService(service string) -> ([]PnPDevice, error)

type PnPDevice struct {
    DeviceID    string
    Name        string
    Service     string
    Status      string
    ClassGuid   string
    Manufacturer string
}
```

---

### 6. Registry Monitoring (🆕 New Package: registry/)

**registry/monitor.go** - Registry Change Notifications
```go
🆕 NewMonitor(keyPath string) -> *Monitor
🆕 Start(callback ChangeCallback) -> error
🆕 Stop() -> error
🆕 WatchKey(key, subkey string, recursive bool) -> error

type ChangeCallback func(change *RegistryChange)

type RegistryChange struct {
    KeyPath   string
    ValueName string
    Operation string // "created", "modified", "deleted"
    OldValue  interface{}
    NewValue  interface{}
}
```

**registry/driver_keys.go** - Driver Registry Parsing
```go
🆕 GetDriverKeys() -> ([]string, error)
🆕 GetDriverConfig(serviceName string) -> (*DriverConfig, error)
🆕 GetDriverParameters(serviceName string) -> (map[string]interface{}, error)
🆕 MonitorDriverChanges(callback) -> error

type DriverConfig struct {
    ImagePath    string
    DisplayName  string
    Start        uint32
    Type         uint32
    ErrorControl uint32
    Group        string
}
```

**registry/device_keys.go** - Device Registry Enumeration
```go
🆕 EnumerateDeviceClasses() -> ([]GUID, error)
🆕 GetDeviceInstanceID(devicePath string) -> (string, error)
🆕 GetDeviceProperties(instanceID string) -> (map[string]interface{}, error)
```

---

### 7. Memory Manipulation (🆕 New Package: memory/)

**memory/process_memory.go** - Process Memory Operations
```go
🆕 ReadMemory(pid uint32, addr uintptr, size uintptr) -> ([]byte, error)
🆕 WriteMemory(pid uint32, addr uintptr, data []byte) -> error
🆕 ReadMemoryEx(process Handle, addr, size uintptr) -> ([]byte, error)
🆕 WriteMemoryEx(process Handle, addr uintptr, data []byte) -> error
🆕 SearchMemory(pid uint32, pattern []byte) -> ([]uintptr, error)
```

**memory/protection.go** - Memory Protection
```go
🆕 ChangeProtection(addr, size uintptr, newProtect uint32) -> (uint32, error)
🆕 ChangeProtectionEx(process Handle, addr, size uintptr, protect uint32) -> (uint32, error)
🆕 MakeExecutable(addr, size uintptr) -> error
🆕 MakeWritable(addr, size uintptr) -> error
```

**memory/allocation.go** - Memory Allocation
```go
🆕 Allocate(size uintptr, protect uint32) -> (uintptr, error)
🆕 AllocateEx(process Handle, addr, size uintptr, protect uint32) -> (uintptr, error)
🆕 Free(addr uintptr) -> error
🆕 FreeEx(process Handle, addr uintptr) -> error
🆕 AllocateExecutable(size uintptr) -> (uintptr, error)
```

---

### 8. PE Parsing (🆕 New Package: pe/)

**pe/parser.go** - PE File Parser
```go
🆕 ParsePE(data []byte) -> (*PEFile, error)
🆕 ParsePEFromFile(path string) -> (*PEFile, error)
🆕 ParsePEFromMemory(baseAddr uintptr) -> (*PEFile, error)
🆕 GetDOSHeader(data []byte) -> (*IMAGE_DOS_HEADER, error)
🆕 GetNTHeaders(data []byte) -> (*IMAGE_NT_HEADERS, error)

type PEFile struct {
    DOSHeader    *IMAGE_DOS_HEADER
    NTHeaders    *IMAGE_NT_HEADERS
    Sections     []*IMAGE_SECTION_HEADER
    Imports      []*ImportDescriptor
    Exports      *ExportDirectory
    BaseAddress  uintptr
}
```

**pe/imports.go** - Import Table Parsing
```go
🆕 ParseImports(pe *PEFile) -> ([]*ImportDescriptor, error)
🆕 GetIATEntry(pe *PEFile, moduleName, funcName string) -> (uintptr, error)
🆕 GetImportedFunctions(pe *PEFile, moduleName string) -> ([]string, error)

type ImportDescriptor struct {
    ModuleName string
    Functions  []ImportFunction
}

type ImportFunction struct {
    Name    string
    Ordinal uint16
    Address uintptr
}
```

**pe/exports.go** - Export Table Parsing
```go
🆕 ParseExports(pe *PEFile) -> (*ExportDirectory, error)
🆕 GetExportByName(pe *PEFile, name string) -> (uintptr, error)
🆕 GetExportByOrdinal(pe *PEFile, ordinal uint16) -> (uintptr, error)

type ExportDirectory struct {
    ModuleName string
    Functions  []ExportFunction
}

type ExportFunction struct {
    Name    string
    Ordinal uint16
    RVA     uint32
    Address uintptr
}
```

---

### 9. IPC Framework (🆕 New Package: ipc/)

**ipc/named_pipe.go** - Named Pipe Communication
```go
🆕 CreatePipeServer(pipeName string) -> (*PipeServer, error)
🆕 ConnectPipeClient(pipeName string) -> (*PipeClient, error)
🆕 Read(buffer []byte) -> (int, error)
🆕 Write(data []byte) -> (int, error)
🆕 Close() -> error

type PipeServer struct {
    PipeName string
    Handle   Handle
    clients  []*PipeClient
}

type PipeClient struct {
    PipeName string
    Handle   Handle
}
```

**ipc/shared_memory.go** - Shared Memory Sections
```go
🆕 CreateSharedMemory(name string, size uintptr) -> (*SharedMemory, error)
🆕 OpenSharedMemory(name string) -> (*SharedMemory, error)
🆕 Write(offset uintptr, data []byte) -> error
🆕 Read(offset, length uintptr) -> ([]byte, error)
🆕 Close() -> error

type SharedMemory struct {
    Name    string
    Handle  Handle
    Mapping uintptr
    Size    uintptr
}
```

**ipc/mailslot.go** - Mailslot Communication
```go
🆕 CreateMailslot(name string) -> (*Mailslot, error)
🆕 OpenMailslot(name string) -> (*Mailslot, error)
🆕 SendMessage(data []byte) -> error
🆕 ReceiveMessage() -> ([]byte, error)

type Mailslot struct {
    Name   string
    Handle Handle
}
```

---

### 10. Assembly Framework (🆕 New Package: asm/)

**asm/x64_encoder.go** - x64 Instruction Encoding
```go
🆕 EncodeInstruction(mnemonic string, operands ...Operand) -> ([]byte, error)
🆕 EncodeJmp(target uintptr, is32bit bool) -> []byte
🆕 EncodeCall(target uintptr) -> []byte
🆕 EncodePush(value uint64) -> []byte
🆕 EncodePop(register Register) -> []byte
🆕 EncodeRet() -> []byte
🆕 EncodeNop(count int) -> []byte

type Register uint8
const (
    RAX Register = iota
    RCX
    RDX
    RBX
    RSP
    RBP
    RSI
    RDI
    R8
    R9
    R10
    R11
    R12
    R13
    R14
    R15
)
```

**asm/jump_gen.go** - Jump Generation
```go
🆕 GenerateAbsoluteJump(target uintptr) -> []byte
🆕 GenerateRelativeJump(source, target uintptr) -> []byte
🆕 GenerateConditionalJump(condition JumpCondition, target uintptr) -> []byte
🆕 CalculateJumpOffset(source, target uintptr) -> int32

type JumpCondition uint8
const (
    JE JumpCondition = iota  // Jump if equal
    JNE                      // Jump if not equal
    JG                       // Jump if greater
    JL                       // Jump if less
    // ... more conditions
)
```

**asm/disasm.go** - Basic Disassembler
```go
🆕 DisassembleInstruction(addr uintptr) -> (*Instruction, error)
🆕 GetInstructionLength(addr uintptr) -> (int, error)
🆕 DisassembleBytes(data []byte, count int) -> ([]*Instruction, error)

type Instruction struct {
    Address  uintptr
    Bytes    []byte
    Mnemonic string
    Operands []string
    Length   int
}
```

---

### 11. Capture System (🆕 New Package: capture/)

**capture/session.go** - Capture Session Management
```go
🆕 NewSession(outputDir string) -> *Session
🆕 Start() -> error
🆕 Stop() -> error
🆕 CaptureIOCTL(device Handle, code uint32, in, out []byte) -> error
🆕 GetCaptureCount() -> uint32
🆕 FlushToDisk() -> error

type Session struct {
    OutputDir    string
    CaptureCount uint32
    Captures     []*Capture
    IsRecording  bool
}

type Capture struct {
    ID         uint32
    Timestamp  time.Time
    DevicePath string
    IOCTLCode  uint32
    InputData  []byte
    OutputData []byte
}
```

**capture/file_format.go** - .conf and .data File I/O
```go
🆕 WriteConfigFile(capture *Capture, filepath string) -> error
🆕 ReadConfigFile(filepath string) -> (*CaptureConfig, error)
🆕 WriteDataFile(capture *Capture, filepath string) -> error
🆕 ReadDataFile(filepath string) -> (*CaptureData, error)
🆕 ParseConfigFormat(data []byte) -> (*CaptureConfig, error)
🆕 SerializeConfig(config *CaptureConfig) -> []byte
```

**capture/buffer_log.go** - Buffer Logging
```go
🆕 LogBuffer(bufferType string, data []byte) -> string
🆕 FormatHexDump(data []byte) -> string
🆕 ParseHexDump(hexStr string) -> ([]byte, error)
🆕 CompareBuffers(buf1, buf2 []byte) -> *BufferDiff

type BufferDiff struct {
    Equal      bool
    Differences []DiffRange
}

type DiffRange struct {
    Offset int
    Length int
    Old    []byte
    New    []byte
}
```

**capture/hook_bridge.go** - Bridge Hooks to Capture
```go
🆕 RegisterHookCallback(callback HookCallback) -> error
🆕 UnregisterHookCallback() -> error
🆕 OnDeviceIoControlCalled(params *IOCTLParams) -> error
🆕 ForwardToSession(session *Session, params *IOCTLParams) -> error

type HookCallback func(params *IOCTLParams) error

type IOCTLParams struct {
    Device      Handle
    IOCTLCode   uint32
    InputBuffer []byte
    OutputBuffer []byte
    BytesReturned uint32
}
```

---

### 12. USB Storage Parser (🆕 Extended in device/)

**device/usb.go** - USB Device History Parser
```go
🆕 GetUSBDeviceHistory() -> ([]USBDevice, error)
🆕 GetUSBDeviceBySerial(serial string) -> (*USBDevice, error)
🆕 GetUSBConnectionTimeline(serial string) -> ([]USBConnection, error)
🆕 EnumerateUSBStorageDevices() -> ([]USBStorageDevice, error)
🆕 ParseSetupAPILog(logPath string) -> ([]USBEvent, error)
🆕 GetUSBSTORRegistryEntries() -> ([]USBSTOREntry, error)

type USBDevice struct {
    SerialNumber    string
    VendorID        uint16
    ProductID       uint16
    FriendlyName    string
    DeviceClass     string
    FirstConnected  time.Time
    LastConnected   time.Time
    DriveLetter     string
    VolumeGUID      string
    ContainerID     string
}

type USBConnection struct {
    Timestamp       time.Time
    EventType       string  // "connect", "disconnect", "mount", "unmount"
    SerialNumber    string
    DriveLetter     string
}

type USBSTOREntry struct {
    DeviceType      string
    VendorProduct   string
    Version         string
    SerialNumber    string
}
```

---

### 13. Windows Artifact Parsers (🆕 New Package: internal/artifacts/)

**internal/artifacts/prefetch/parser.go** - Windows Prefetch Parser (pf)
```go
🆕 ParsePrefetchFile(path string) -> (*PrefetchFile, error)
🆕 ParsePrefetchDir(dirPath string) -> ([]*PrefetchFile, error)
🆕 GetExecutionTimeline(prefetchFiles []*PrefetchFile) -> ([]ExecutionEvent, error)
🆕 DecompressMAM(data []byte) -> ([]byte, error)

type PrefetchFile struct {
    Version         uint32    // 17, 23, 26, 30
    Signature       uint32
    FileSize        uint32
    ExecutableName  string
    PrefetchHash    uint32
    RunCount        uint32
    LastRunTimes    []time.Time  // Up to 8 timestamps
    FilesReferenced []FileReference
    VolumesInfo     []VolumeInfo
}

type FileReference struct {
    Filename        string
    NTFSReference   uint64
}

type ExecutionEvent struct {
    Executable      string
    Timestamp       time.Time
    RunCount        uint32
    PrefetchPath    string
}
```

**internal/artifacts/lnk/parser.go** - Windows LNK Parser (lp)
```go
🆕 ParseLNKFile(path string) -> (*ShellLink, error)
🆕 ParseLNKBytes(data []byte) -> (*ShellLink, error)
🆕 ResolveTarget(lnk *ShellLink) -> (string, error)
🆕 GetExtraDataBlocks(lnk *ShellLink) -> ([]ExtraDataBlock, error)

type ShellLink struct {
    HeaderSize          uint32
    LinkFlags           uint32
    FileAttributes      uint32
    CreationTime        time.Time
    AccessTime          time.Time
    WriteTime           time.Time
    FileSize            uint32
    IconIndex           int32
    ShowCommand         uint32
    HotKey              uint16
    TargetPath          string
    Arguments           string
    WorkingDir          string
    IconLocation        string
    Description         string
    LinkInfo            *LinkInfo
    ExtraData           []ExtraDataBlock
}

type LinkInfo struct {
    VolumeID            *VolumeID
    LocalBasePath       string
    CommonNetworkRelativeLink *CommonNetworkRelativeLink
    CommonPathSuffix    string
}
```

**internal/artifacts/jumplist/parser.go** - Windows Jump List Parser (jmp)
```go
🆕 ParseAutomaticDestinations(path string) -> (*JumpList, error)
🆕 ParseCustomDestinations(path string) -> (*JumpList, error)
🆕 ParseDestListStream(data []byte) -> ([]DestListEntry, error)
🆕 ExtractLNKEntries(jumpList *JumpList) -> ([]*ShellLink, error)

type JumpList struct {
    AppID               string
    Type                string  // "automatic" or "custom"
    Entries             []JumpListEntry
    DestListHeader      *DestListHeader
}

type JumpListEntry struct {
    EntryID             uint32
    ShellLink           *ShellLink
    Hostname            string
    NetBIOSName         string
    EntryIDLow          uint32
    BirthDate           time.Time
    MACAddress          string
    AccessCount         uint32
    LastAccessTime      time.Time
    PinStatus           int32
}
```

**internal/artifacts/shellbag/parser.go** - Windows ShellBag Parser (sbag)
```go
🆕 ParseShellBags(registryPath string) -> (*ShellBagTree, error)
🆕 ParseShellBagsFromHive(hive *RegistryHive) -> (*ShellBagTree, error)
🆕 ParseItemIDList(data []byte) -> (*ItemIDList, error)
🆕 ExtractShellItems(itemIDList *ItemIDList) -> ([]ShellItem, error)

type ShellBagTree struct {
    Root            *ShellBagNode
    TotalEntries    int
}

type ShellBagNode struct {
    Path            string
    ShellItem       ShellItem
    LastWriteTime   time.Time
    Children        []*ShellBagNode
}

type ShellItem struct {
    Type            uint8
    Size            uint16
    Name            string
    ModifiedTime    time.Time
    CreatedTime     time.Time
    AccessedTime    time.Time
    MFTEntryNumber  uint64
    MFTSequenceNum  uint16
}
```

**internal/artifacts/shimcache/parser.go** - AppCompatibility Cache (wacu)
```go
🆕 ParseShimCache(registryPath string) -> (*ShimCache, error)
🆕 ParseShimCacheFromHive(hive *RegistryHive) -> (*ShimCache, error)
🆕 ParseShimCacheWin10(data []byte) -> ([]ShimCacheEntry, error)
🆕 ParseShimCacheWin8(data []byte) -> ([]ShimCacheEntry, error)
🆕 ParseShimCacheWin7(data []byte) -> ([]ShimCacheEntry, error)

type ShimCache struct {
    Version         string
    Entries         []ShimCacheEntry
}

type ShimCacheEntry struct {
    Path            string
    LastModified    time.Time
    Size            uint64
    ExecFlag        bool
    InsertedTime    time.Time
    ShimFlags       uint32
    DataSize        uint32
}
```

**internal/artifacts/shimdb/parser.go** - Shim Database Parser (shims)
```go
🆕 ParseSDBFile(path string) -> (*SDBDatabase, error)
🆕 GetDatabaseInfo(sdb *SDBDatabase) -> (*SDBInfo, error)
🆕 EnumerateFixes(sdb *SDBDatabase) -> ([]SDBFix, error)
🆕 GetApplicationShims(appName string) -> ([]SDBShim, error)

type SDBDatabase struct {
    Magic           uint32
    MajorVersion    uint16
    MinorVersion    uint16
    Tags            []SDBTag
    StringTable     map[uint32]string
    Indexes         []SDBIndex
}

type SDBTag struct {
    Type            uint16
    TagID           uint16
    Size            uint32
    Data            interface{}
    Children        []SDBTag
}

type SDBFix struct {
    Name            string
    AppName         string
    Vendor          string
    ExePath         string
    Shims           []SDBShim
    Flags           []SDBFlag
    MatchMode       string
}
```

**internal/artifacts/activitiescache/parser.go** - ActivitiesCache Parser (tac)
```go
🆕 ParseActivitiesCache(dbPath string) -> (*ActivitiesCache, error)
🆕 GetActivityTimeline(cache *ActivitiesCache) -> ([]Activity, error)
🆕 FilterActivitiesByType(activities []Activity, actType string) -> []Activity
🆕 FilterActivitiesByDateRange(activities []Activity, start, end time.Time) -> []Activity

type ActivitiesCache struct {
    Activities      []Activity
    ActivityOperations []ActivityOperation
    SmartLookup     []SmartLookupEntry
}

type Activity struct {
    ID              string
    AppID           string
    PackageIDHash   string
    AppActivityID   string
    ActivityType    int32
    ActivityStatus  int32
    ParentActivityID string
    Tag             string
    Group           string
    MatchID         string
    LastModifiedTime time.Time
    ExpirationTime  time.Time
    Payload         string
    Priority        int32
    StartTime       time.Time
    EndTime         time.Time
    Duration        int64
    Platform        string
}
```

**internal/artifacts/recycle/parser.go** - Recycle Bin Parser (tia)
```go
🆕 ParseRecycleBin(recyclePath string) -> (*RecycleBin, error)
🆕 ParseINFO2(path string) -> ([]DeletedFile, error)      // Windows XP
🆕 ParseIDollarFile(path string) -> (*DeletedFile, error) // Vista+
🆕 GetDeletedFileTimeline(bin *RecycleBin) -> ([]DeletedFile, error)
🆕 IdentifyRecoverableFiles(bin *RecycleBin) -> ([]DeletedFile, error)

type RecycleBin struct {
    Version         string
    DeletedFiles    []DeletedFile
    Location        string
    SID             string
}

type DeletedFile struct {
    OriginalPath    string
    DeletedPath     string
    DeletionTime    time.Time
    FileSize        int64
    IFileName       string  // $I file name
    RFileName       string  // $R file name
    IsDirectory     bool
    Recoverable     bool
}
```

**internal/artifacts/wpn/parser.go** - Windows Push Notification Parser (wpn)
```go
🆕 ParseWPNDatabase(dbPath string) -> (*WPNDatabase, error)
🆕 GetNotifications(db *WPNDatabase) -> ([]Notification, error)
🆕 GetNotificationHandlers(db *WPNDatabase) -> ([]NotificationHandler, error)

type WPNDatabase struct {
    Notifications       []Notification
    Handlers            []NotificationHandler
    Settings            []NotificationSetting
}

type Notification struct {
    Order               int64
    HandlerID           int64
    WNSCreatedTime      time.Time
    ExpiryTime          time.Time
    ArrivalTime         time.Time
    PayloadType         string
    Payload             string
    Tag                 string
    Group               string
}

type NotificationHandler struct {
    ID                  int64
    PrimaryID           string
    WNSChannelID        string
    WNSChannelSecondary string
    CreatedTime         time.Time
    ModifiedTime        time.Time
    AppID               string
}
```

**internal/artifacts/backstage/parser.go** - MS Office Backstage Parser (bs)
```go
🆕 ParseOfficeMRU(registryPath string) -> (*OfficeMRU, error)
🆕 ParseOfficeMRUFromHive(hive *RegistryHive) -> (*OfficeMRU, error)
🆕 GetRecentDocuments(mru *OfficeMRU) -> ([]OfficeDocument, error)
🆕 GetOfficeVersion() -> (string, error)

type OfficeMRU struct {
    Version         string
    Applications    map[string]*ApplicationMRU  // Word, Excel, PowerPoint, etc.
}

type ApplicationMRU struct {
    AppName         string
    RecentFiles     []OfficeDocument
    PlaceMRU        []OfficePlace
    TrustRecords    []TrustRecord
}

type OfficeDocument struct {
    Path            string
    LastAccess      time.Time
    Position        int
    ReadOnly        bool
}
```

---

### 14. Browser Artifact Parsers (🆕 New Package: internal/browser/)

**internal/browser/chromium/parser.go** - Chromium Parser (csp)
```go
🆕 NewChromiumParser(profilePath string) -> (*ChromiumParser, error)
🆕 GetHistory() -> ([]HistoryEntry, error)
🆕 GetCookies() -> ([]Cookie, error)
🆕 GetDownloads() -> ([]Download, error)
🆕 GetBookmarks() -> ([]Bookmark, error)
🆕 GetLogins() -> ([]Login, error)  // Encrypted data
🆕 GetAutofill() -> ([]AutofillEntry, error)
🆕 GetExtensions() -> ([]Extension, error)

type ChromiumParser struct {
    ProfilePath     string
    Browser         string  // Chrome, Edge, Brave, etc.
}

type HistoryEntry struct {
    URL             string
    Title           string
    VisitCount      int
    TypedCount      int
    LastVisit       time.Time
    Hidden          bool
    Transition      int
}

type Cookie struct {
    Host            string
    Name            string
    Value           string
    Path            string
    ExpiresUTC      time.Time
    Secure          bool
    HTTPOnly        bool
    SameSite        int
}

type Download struct {
    TargetPath      string
    URL             string
    StartTime       time.Time
    EndTime         time.Time
    TotalBytes      int64
    ReceivedBytes   int64
    State           int
    InterruptReason int
    MimeType        string
}
```

**internal/browser/chromium/cache/parser.go** - Chromium Cache Parser (ccp)
```go
🆕 ParseCacheIndex(indexPath string) -> (*CacheIndex, error)
🆕 GetCacheEntries(cachePath string) -> ([]CacheEntry, error)
🆕 ExtractCacheEntry(entry *CacheEntry, outputPath string) -> error
🆕 SearchCacheByURL(pattern string) -> ([]CacheEntry, error)

type CacheIndex struct {
    Version         uint32
    EntryCount      uint32
    AddressTable    []uint32
}

type CacheEntry struct {
    Hash            uint32
    URL             string
    CreationTime    time.Time
    AccessTime      time.Time
    DataSize        int64
    ContentType     string
    CacheAddress    uint32
    DataFile        int
    DataOffset      int
}
```

**internal/browser/mozilla/parser.go** - Mozilla Parser (msp)
```go
🆕 NewMozillaParser(profilePath string) -> (*MozillaParser, error)
🆕 GetPlaces() -> (*Places, error)  // History + Bookmarks
🆕 GetHistory() -> ([]HistoryEntry, error)
🆕 GetBookmarks() -> ([]Bookmark, error)
🆕 GetCookies() -> ([]Cookie, error)
🆕 GetFormHistory() -> ([]FormEntry, error)
🆕 GetLogins() -> ([]Login, error)
🆕 GetDownloads() -> ([]Download, error)

type MozillaParser struct {
    ProfilePath     string
    Browser         string  // Firefox, Thunderbird, etc.
}

type Places struct {
    History         []HistoryEntry
    Bookmarks       []Bookmark
    Annotations     []Annotation
    InputHistory    []InputHistoryEntry
}

type FormEntry struct {
    FieldName       string
    Value           string
    TimesUsed       int
    FirstUsed       time.Time
    LastUsed        time.Time
}
```

**internal/browser/mozilla/cache/parser.go** - Mozilla Cache Parser (mcp)
```go
🆕 ParseCache2Index(indexPath string) -> (*Cache2Index, error)
🆕 GetCacheEntries(cachePath string) -> ([]Cache2Entry, error)
🆕 ExtractCacheEntry(entry *Cache2Entry, outputPath string) -> error
🆕 ParseCacheMetadata(entry *Cache2Entry) -> (*CacheMetadata, error)

type Cache2Index struct {
    Version         uint32
    LastClean       time.Time
    DirtyFlag       bool
}

type Cache2Entry struct {
    Hash            string
    URL             string
    FetchCount      int32
    LastFetch       time.Time
    LastModified    time.Time
    ExpirationTime  time.Time
    DataSize        int64
    MetadataSize    int32
    ContentType     string
}
```

**internal/browser/safari/parser.go** - Safari Parser (sap)
```go
🆕 NewSafariParser(profilePath string) -> (*SafariParser, error)
🆕 GetHistory() -> ([]HistoryEntry, error)
🆕 GetDownloads() -> ([]Download, error)
🆕 GetBookmarks() -> ([]Bookmark, error)
🆕 GetRecentSearches() -> ([]SearchEntry, error)

type SafariParser struct {
    ProfilePath     string
}

type SearchEntry struct {
    SearchDescriptor string
    SearchQuery     string
    Timestamp       time.Time
}
```

---

### 15. Registry Analysis Extended (🆕 Extended in registry/)

**registry/hive.go** - Registry Hive Parser (yaru)
```go
🆕 OpenHive(path string) -> (*RegistryHive, error)
🆕 OpenHiveWithTransactions(path string, logFiles []string) -> (*RegistryHive, error)
🆕 EnumerateKeys(hive *RegistryHive, keyPath string) -> ([]RegistryKey, error)
🆕 GetValue(hive *RegistryHive, keyPath, valueName string) -> (*RegistryValue, error)
🆕 GetAllValues(hive *RegistryHive, keyPath string) -> ([]RegistryValue, error)
🆕 SearchKeys(hive *RegistryHive, pattern string) -> ([]RegistryKey, error)
🆕 SearchValues(hive *RegistryHive, pattern string) -> ([]RegistryValue, error)
🆕 GetDeletedKeys(hive *RegistryHive) -> ([]DeletedKey, error)

type RegistryHive struct {
    Path            string
    RootKey         *RegistryKey
    BaseBlock       *HiveBaseBlock
    IsDirty         bool
}

type HiveBaseBlock struct {
    Signature       uint32  // "regf"
    PrimarySequence uint32
    SecondarySequence uint32
    LastWriteTime   time.Time
    MajorVersion    uint32
    MinorVersion    uint32
    Type            uint32
    Format          uint32
    RootCellOffset  uint32
    HiveBinsSize    uint32
}

type RegistryKey struct {
    Name            string
    Path            string
    ClassName       string
    LastWriteTime   time.Time
    SubKeyCount     uint32
    ValueCount      uint32
    SecurityDescriptor []byte
    Flags           uint16
}

type RegistryValue struct {
    Name            string
    Type            uint32  // REG_SZ, REG_BINARY, etc.
    Data            interface{}
    DataSize        uint32
}

type DeletedKey struct {
    Key             RegistryKey
    DeletionTime    time.Time  // Approximate
    RecoveredFrom   string
}
```

---

### 16. Event Log Analysis (🆕 New Package: internal/evtx/)

**internal/evtx/parser.go** - EVTX Parser (evtwalk/evtx_view)
```go
🆕 OpenEVTX(path string) -> (*EVTXFile, error)
🆕 GetAllRecords() -> ([]EventRecord, error)
🆕 GetRecordsByEventID(eventID uint16) -> ([]EventRecord, error)
🆕 GetRecordsByTimeRange(start, end time.Time) -> ([]EventRecord, error)
🆕 GetRecordsByProvider(providerName string) -> ([]EventRecord, error)
🆕 SearchRecords(query string) -> ([]EventRecord, error)
🆕 ExportToXML(records []EventRecord, outputPath string) -> error
🆕 ExportToJSON(records []EventRecord, outputPath string) -> error

type EVTXFile struct {
    Path            string
    Header          *EVTXFileHeader
    Chunks          []*EVTXChunk
    RecordCount     uint64
}

type EVTXFileHeader struct {
    Magic           [8]byte  // "ElfFile\x00"
    FirstChunkNum   uint64
    LastChunkNum    uint64
    NextRecordID    uint64
    HeaderSize      uint32
    MinorVersion    uint16
    MajorVersion    uint16
    ChunkCount      uint16
    Flags           uint32
}

type EventRecord struct {
    RecordID        uint64
    Timestamp       time.Time
    EventID         uint16
    Level           uint8
    Channel         string
    Provider        string
    Computer        string
    UserSID         string
    ProcessID       uint32
    ThreadID        uint32
    Keywords        uint64
    EventData       map[string]interface{}
    XMLData         string
}
```

**internal/evtx/fragment/parser.go** - EVTX Fragment Parser (evtfx)
```go
🆕 CarveEVTXFragments(imagePath string) -> ([]EVTXFragment, error)
🆕 RecoverFragmentedRecords(fragments []EVTXFragment) -> ([]EventRecord, error)
🆕 ValidateFragment(fragment *EVTXFragment) -> bool

type EVTXFragment struct {
    Offset          int64
    Size            int32
    ChunkHeader     *EVTXChunkHeader
    Records         []EventRecord
    IsComplete      bool
}
```

**internal/messagetable/parser.go** - Message Table Parser (elmo)
```go
🆕 ExtractMessageTable(pePath string) -> (*MessageTable, error)
🆕 GetMessage(tableID uint32, messageID uint32) -> (string, error)
🆕 FormatMessage(msg string, insertions []string) -> string
🆕 ListMessageTables(pePath string) -> ([]MessageTableInfo, error)

type MessageTable struct {
    ResourceID      uint32
    Language        uint32
    Messages        map[uint32]string
}

type MessageTableInfo struct {
    ResourceID      uint32
    Language        uint32
    MessageCount    int
}
```

---

### 17. NTFS Filesystem Analysis (🆕 New Package: internal/filesystem/)

**internal/filesystem/ntfs/volume.go** - NTFS Core
```go
🆕 OpenNTFSVolume(path string) -> (*NTFSVolume, error)
🆕 OpenNTFSImage(imagePath string) -> (*NTFSVolume, error)
🆕 GetVolumeInfo() -> (*VolumeInfo, error)
🆕 ReadClusters(lcn uint64, count uint32) -> ([]byte, error)
🆕 ClusterToOffset(lcn uint64) -> int64
🆕 OffsetToCluster(offset int64) -> uint64

type NTFSVolume struct {
    Handle          Handle
    BootSector      *NTFSBootSector
    MFT             *MFTParser
    ClusterSize     uint32
    SectorSize      uint16
    MFTOffset       int64
}

type NTFSBootSector struct {
    OEMIdentifier   string
    BytesPerSector  uint16
    SectorsPerCluster uint8
    MediaDescriptor uint8
    TotalSectors    uint64
    MFTCluster      uint64
    MFTMirrorCluster uint64
    ClustersPerMFTRecord int8
    ClustersPerIndexRecord int8
    VolumeSerialNumber uint64
}

type VolumeInfo struct {
    VolumeName      string
    VolumeVersion   string
    Flags           uint16
    SerialNumber    uint64
    TotalSize       int64
    FreeSpace       int64
}
```

**internal/filesystem/mft/parser.go** - $MFT Parser (ntfswalk)
```go
🆕 OpenMFT(volume *NTFSVolume) -> (*MFTParser, error)
🆕 OpenMFTFile(path string) -> (*MFTParser, error)
🆕 GetRecordByNumber(recordNum uint64) -> (*MFTRecord, error)
🆕 GetRecordByPath(path string) -> (*MFTRecord, error)
🆕 EnumerateRecords() -> (<-chan *MFTRecord, error)
🆕 GetFileTimestamps(record *MFTRecord) -> (*Timestamps, error)
🆕 GetDataRuns(record *MFTRecord) -> ([]DataRun, error)
🆕 GetAlternateDataStreams(record *MFTRecord) -> ([]ADSEntry, error)

type MFTParser struct {
    Volume          *NTFSVolume
    RecordSize      uint32
    RecordCount     uint64
}

type MFTRecord struct {
    RecordNumber    uint64
    SequenceNumber  uint16
    Flags           uint16
    LogFileSeqNum   uint64
    BaseRecordRef   uint64
    Attributes      []MFTAttribute
    IsDeleted       bool
    IsDirectory     bool
    ParentDirRef    uint64
}

type MFTAttribute struct {
    TypeCode        uint32
    Name            string
    Flags           uint16
    IsResident      bool
    ResidentData    []byte
    NonResidentRuns []DataRun
}

type DataRun struct {
    VCN             uint64  // Virtual Cluster Number
    LCN             uint64  // Logical Cluster Number
    Length          uint64  // Length in clusters
}

type Timestamps struct {
    Created         time.Time
    Modified        time.Time
    MFTModified     time.Time
    Accessed        time.Time
    // From $FILE_NAME attribute (if different)
    FNCreated       time.Time
    FNModified      time.Time
    FNMFTModified   time.Time
    FNAccessed      time.Time
}

type ADSEntry struct {
    Name            string
    Size            int64
    DataRuns        []DataRun
}
```

**internal/filesystem/usnjrnl/parser.go** - $UsnJrnl Parser (jp)
```go
🆕 OpenUSNJournal(volume *NTFSVolume) -> (*USNJournalParser, error)
🆕 OpenUSNJournalFile(path string) -> (*USNJournalParser, error)
🆕 GetAllRecords() -> (<-chan *USNRecord, error)
🆕 GetRecordsByReason(reasons uint32) -> ([]USNRecord, error)
🆕 GetRecordsByTimeRange(start, end time.Time) -> ([]USNRecord, error)
🆕 GetRecordsByFilename(pattern string) -> ([]USNRecord, error)

type USNJournalParser struct {
    Volume          *NTFSVolume
    JournalData     []byte
    MaxUSN          uint64
    FirstUSN        uint64
}

type USNRecord struct {
    RecordLength    uint32
    MajorVersion    uint16
    MinorVersion    uint16
    FileReference   uint64
    ParentReference uint64
    USN             uint64
    Timestamp       time.Time
    Reason          uint32
    SourceInfo      uint32
    SecurityID      uint32
    FileAttributes  uint32
    FileName        string
}

// Reason flags
const (
    USN_REASON_DATA_OVERWRITE      = 0x00000001
    USN_REASON_DATA_EXTEND         = 0x00000002
    USN_REASON_DATA_TRUNCATION     = 0x00000004
    USN_REASON_NAMED_DATA_OVERWRITE = 0x00000010
    USN_REASON_NAMED_DATA_EXTEND   = 0x00000020
    USN_REASON_NAMED_DATA_TRUNCATION = 0x00000040
    USN_REASON_FILE_CREATE         = 0x00000100
    USN_REASON_FILE_DELETE         = 0x00000200
    USN_REASON_RENAME_OLD_NAME     = 0x00001000
    USN_REASON_RENAME_NEW_NAME     = 0x00002000
    USN_REASON_SECURITY_CHANGE     = 0x00000800
    USN_REASON_CLOSE               = 0x80000000
)
```

**internal/filesystem/logfile/parser.go** - $LogFile Parser (mala)
```go
🆕 OpenLogFile(volume *NTFSVolume) -> (*LogFileParser, error)
🆕 OpenLogFileFile(path string) -> (*LogFileParser, error)
🆕 GetRestartArea() -> (*RestartArea, error)
🆕 GetLogRecords() -> (<-chan *LogRecord, error)
🆕 AnalyzeTransactions() -> ([]Transaction, error)

type LogFileParser struct {
    Volume          *NTFSVolume
    RestartAreas    [2]*RestartArea
    ClientRecords   []LogRecord
}

type RestartArea struct {
    MajorVersion    uint16
    MinorVersion    uint16
    StartOfCheckpoint uint64
    OpenAttributeTableLSN uint64
    AttributeNamesLSN uint64
    DirtyPageTableLSN uint64
    TransactionTableLSN uint64
}

type LogRecord struct {
    ThisLSN         uint64
    ClientPreviousLSN uint64
    ClientUndoNextLSN uint64
    ClientDataLength uint32
    ClientID        uint32
    RecordType      uint32
    TransactionID   uint32
    Flags           uint16
    RedoOperation   uint16
    UndoOperation   uint16
    RedoData        []byte
    UndoData        []byte
}
```

**internal/filesystem/indx/parser.go** - INDX Slack Parser (wisp)
```go
🆕 OpenINDX(volume *NTFSVolume, directoryRef uint64) -> (*INDXParser, error)
🆕 ParseINDXBuffer(data []byte) -> (*INDXBuffer, error)
🆕 GetSlackEntries() -> ([]SlackEntry, error)
🆕 CarveDeletedEntries() -> ([]DeletedIndexEntry, error)

type INDXParser struct {
    Volume          *NTFSVolume
    DirectoryRef    uint64
    INDXBuffers     []INDXBuffer
}

type INDXBuffer struct {
    Magic           uint32  // "INDX"
    UpdateSeqOffset uint16
    UpdateSeqSize   uint16
    LogFileSeqNum   uint64
    VCN             uint64
    Entries         []IndexEntry
    SlackSpace      []byte
}

type IndexEntry struct {
    FileReference   uint64
    EntryLength     uint16
    ContentLength   uint16
    Flags           uint32
    FileName        string
    FileSize        uint64
    CreatedTime     time.Time
    ModifiedTime    time.Time
}

type SlackEntry struct {
    Offset          int64
    Entry           IndexEntry
    Confidence      float32
    IsRecoverable   bool
}
```

**internal/filesystem/fat/parser.go** - FAT32/exFAT Parser (fata)
```go
🆕 OpenFATVolume(path string) -> (*FATVolume, error)
🆕 OpenFATImage(imagePath string) -> (*FATVolume, error)
🆕 GetVolumeInfo() -> (*FATVolumeInfo, error)
🆕 EnumerateFiles(dirCluster uint32) -> ([]FATDirectoryEntry, error)
🆕 GetDeletedFiles() -> ([]DeletedFATFile, error)
🆕 RecoverFile(entry *DeletedFATFile, outputPath string) -> error

type FATVolume struct {
    Handle          Handle
    Type            string  // "FAT32" or "exFAT"
    BootSector      *FATBootSector
    FAT             []uint32
    ClusterSize     uint32
}

type FATBootSector struct {
    OEMName         string
    BytesPerSector  uint16
    SectorsPerCluster uint8
    ReservedSectors uint16
    NumberOfFATs    uint8
    TotalSectors    uint32
    RootCluster     uint32  // FAT32 only
    VolumeLabel     string
}

type FATDirectoryEntry struct {
    Name            string
    Extension       string
    Attributes      uint8
    CreatedTime     time.Time
    ModifiedTime    time.Time
    AccessedDate    time.Time
    FirstCluster    uint32
    FileSize        uint32
    IsDeleted       bool
    LongFileName    string
}
```

---

### 18. Network Utilities (🆕 New Package: internal/network/)

**internal/network/dns/query.go** - DNS Query Utility (dqu)
```go
🆕 Query(domain string, recordType string) -> ([]DNSRecord, error)
🆕 QueryWithServer(domain, recordType, server string) -> ([]DNSRecord, error)
🆕 GetSystemDNSCache() -> ([]CachedDNSEntry, error)
🆕 ClearDNSCache() -> error
🆕 ReverseLookup(ip string) -> ([]string, error)

type DNSRecord struct {
    Name            string
    Type            string
    TTL             uint32
    Data            string
    Priority        uint16  // MX records
}

type CachedDNSEntry struct {
    Name            string
    Type            string
    Data            string
    TTL             uint32
    ExpirationTime  time.Time
}
```

**internal/network/pcap/parser.go** - PCAP Parser (pic)
```go
🆕 OpenPCAP(path string) -> (*PCAPReader, error)
🆕 OpenPCAPNG(path string) -> (*PCAPNGReader, error)
🆕 GetPackets() -> (<-chan *Packet, error)
🆕 FilterByProtocol(protocol string) -> ([]Packet, error)
🆕 CarveICMPData() -> ([]ICMPPayload, error)
🆕 WritePCAP(packets []Packet, outputPath string) -> error

type PCAPReader struct {
    Path            string
    Header          *PCAPHeader
    LinkType        uint32
}

type PCAPHeader struct {
    MagicNumber     uint32
    VersionMajor    uint16
    VersionMinor    uint16
    ThisZone        int32
    SigFigs         uint32
    SnapLen         uint32
    Network         uint32
}

type Packet struct {
    Timestamp       time.Time
    CapturedLength  uint32
    OriginalLength  uint32
    Data            []byte
    Ethernet        *EthernetHeader
    IP              *IPHeader
    TCP             *TCPHeader
    UDP             *UDPHeader
    ICMP            *ICMPHeader
}

type ICMPPayload struct {
    Timestamp       time.Time
    Type            uint8
    Code            uint8
    SourceIP        string
    DestIP          string
    Payload         []byte
    SequenceNum     uint16
}
```

**internal/network/netxfer/server.go** - Network Transfer (nx)
```go
🆕 NewServer(bindAddr string, port int) -> (*TransferServer, error)
🆕 Start() -> error
🆕 Stop() -> error
🆕 SetEncryption(enabled bool, key []byte) -> error
🆕 OnFileReceived(callback FileReceivedCallback) -> error

type TransferServer struct {
    BindAddr        string
    Port            int
    Encrypted       bool
    Connections     []*Connection
}

type FileReceivedCallback func(filename string, data []byte, metadata *FileMetadata)

type FileMetadata struct {
    Filename        string
    Size            int64
    Checksum        string
    Timestamp       time.Time
    Sender          string
}
```

---

### 19. PE Extended Analysis (🆕 New Package: internal/pe/)

**internal/pe/scanner.go** - PE Anomaly Scanner (pescan)
```go
🆕 ScanPE(path string) -> (*ScanResult, error)
🆕 ScanPEBytes(data []byte) -> (*ScanResult, error)
🆕 DetectPackers() -> ([]PackerSignature, error)
🆕 DetectAnomalies() -> ([]Anomaly, error)
🆕 CalculateEntropy(section *Section) -> float64
🆕 GetImphash() -> string
🆕 ValidateSignature() -> (*SignatureInfo, error)

type ScanResult struct {
    Path            string
    Imphash         string
    Entropy         float64
    IsPacked        bool
    Packers         []PackerSignature
    Anomalies       []Anomaly
    Signature       *SignatureInfo
    Imports         []ImportInfo
    Exports         []ExportInfo
    Resources       []ResourceInfo
}

type Anomaly struct {
    Type            string
    Description     string
    Severity        string  // "low", "medium", "high"
    Location        string
}

type PackerSignature struct {
    Name            string
    Version         string
    Confidence      float32
    Offset          int64
}

type SignatureInfo struct {
    IsSigned        bool
    SignerName      string
    Issuer          string
    Timestamp       time.Time
    IsValid         bool
    ErrorMessage    string
}
```

**internal/pe/viewer.go** - PE Viewer (pe_view)
```go
🆕 GetPEInfo(path string) -> (*PEInfo, error)
🆕 GetHeaders() -> (*Headers, error)
🆕 GetSections() -> ([]SectionInfo, error)
🆕 GetImports() -> ([]ImportDLL, error)
🆕 GetExports() -> ([]ExportFunc, error)
🆕 GetResources() -> ([]ResourceEntry, error)
🆕 GetVersionInfo() -> (*VersionInfo, error)
🆕 GetManifest() -> (string, error)
🆕 GetDebugInfo() -> (*DebugInfo, error)
🆕 DumpSection(sectionName string, outputPath string) -> error

type PEInfo struct {
    Path            string
    FileSize        int64
    MD5             string
    SHA1            string
    SHA256          string
    Type            string  // "EXE", "DLL", "SYS"
    Subsystem       string
    Machine         string
    Timestamp       time.Time
    EntryPoint      uint64
    ImageBase       uint64
    Characteristics uint16
    DLLCharacteristics uint16
}

type VersionInfo struct {
    FileVersion     string
    ProductVersion  string
    CompanyName     string
    FileDescription string
    InternalName    string
    LegalCopyright  string
    OriginalFilename string
    ProductName     string
}
```

---

### 20. Miscellaneous Utilities (🆕 New Packages: internal/vss/, internal/symbol/, etc.)

**internal/vss/enumerator.go** - VSS Enumerator (vssenum)
```go
🆕 EnumerateShadowCopies() -> ([]ShadowCopy, error)
🆕 GetShadowCopyByID(id string) -> (*ShadowCopy, error)
🆕 MountShadowCopy(shadowID string, mountPoint string) -> error
🆕 UnmountShadowCopy(mountPoint string) -> error
🆕 CopyFileFromShadow(shadowID, filePath, outputPath string) -> error
🆕 DiffSnapshots(shadowID1, shadowID2 string) -> (*SnapshotDiff, error)

type ShadowCopy struct {
    ID              string
    SetID           string
    SnapshotTime    time.Time
    OriginalVolume  string
    DeviceObject    string
    State           uint32
    Attributes      uint32
    Provider        string
}

type SnapshotDiff struct {
    AddedFiles      []string
    DeletedFiles    []string
    ModifiedFiles   []DiffFile
}

type DiffFile struct {
    Path            string
    OldSize         int64
    NewSize         int64
    OldHash         string
    NewHash         string
}
```

**internal/symbol/fetch.go** - Symbol Fetch Utility (sf)
```go
🆕 FetchSymbol(modulePath string) -> (*PDBInfo, error)
🆕 FetchSymbolByGUID(guid, age string, pdbName string) -> (string, error)
🆕 SetSymbolServer(url string) -> error
🆕 SetCachePath(path string) -> error
🆕 GetCachedSymbols() -> ([]CachedSymbol, error)
🆕 ExtractGUIDFromPE(pePath string) -> (*DebugGUID, error)

type PDBInfo struct {
    PDBName         string
    GUID            string
    Age             uint32
    LocalPath       string
    RemoteURL       string
}

type DebugGUID struct {
    GUID            string
    Age             uint32
    PDBName         string
}

type CachedSymbol struct {
    PDBName         string
    GUID            string
    LocalPath       string
    DownloadTime    time.Time
    Size            int64
}
```

**internal/csvdx/parser.go** - CSV Utilities (csvdx)
```go
🆕 ReadCSV(path string) -> (*CSVData, error)
🆕 WriteCSV(data *CSVData, path string) -> error
🆕 MergeCSVFiles(paths []string, outputPath string) -> error
🆕 TransformColumn(data *CSVData, column string, transform TransformFunc) -> error
🆕 FilterRows(data *CSVData, predicate FilterFunc) -> *CSVData
🆕 SortByColumn(data *CSVData, column string, ascending bool) -> error

type CSVData struct {
    Headers         []string
    Rows            [][]string
    Delimiter       rune
}

type TransformFunc func(value string) string
type FilterFunc func(row []string) bool
```

**internal/disk/imaging.go** - Disk Utility (dup)
```go
🆕 OpenDisk(path string) -> (*DiskHandle, error)
🆕 ReadSectors(startSector, count uint64) -> ([]byte, error)
🆕 GetPartitions() -> ([]Partition, error)
🆕 ParseMBR() -> (*MBR, error)
🆕 ParseGPT() -> (*GPT, error)
🆕 CreateForensicImage(outputPath string, opts *ImagingOptions) -> error
🆕 VerifyImage(imagePath, originalPath string) -> (bool, error)

type DiskHandle struct {
    Path            string
    Handle          Handle
    Size            int64
    SectorSize      uint32
    Geometry        *DiskGeometry
}

type Partition struct {
    Number          int
    Type            string
    StartSector     uint64
    EndSector       uint64
    Size            int64
    Bootable        bool
    FileSystem      string
    VolumeLabel     string
}

type MBR struct {
    BootCode        []byte
    Partitions      [4]MBRPartition
    Signature       uint16
}

type GPT struct {
    Header          *GPTHeader
    Partitions      []GPTPartition
}

type ImagingOptions struct {
    Compression     string  // "none", "gzip", "zstd"
    Split           bool
    SplitSize       int64
    HashAlgorithm   string  // "md5", "sha1", "sha256"
    Verify          bool
}
```

---

## Complete Workflow Examples

### Workflow 1: Load Driver and Discover IOCTLs

```
1. User calls: service.LoadDriver("C:\\mydriver.sys", "MyDriver")
   ├─> Opens Service Control Manager
   ├─> Creates service entry
   ├─> Starts driver service
   └─> Returns service handle

2. User calls: device.FindDevicesByService("MyDriver")
   ├─> Queries SetupAPI for devices with service="MyDriver"
   ├─> Returns list of DeviceInfo with device paths
   └─> Example: \\?\Device\MyDriver, \\.\MyDevice

3. User calls: device.FuzzIOCTL("\\.\MyDevice", opts)
   ├─> Generates IOCTL codes (0x220000 - 0x220FFF)
   ├─> For each code:
   │   ├─> Opens device handle
   │   ├─> Calls DeviceIoControl with code
   │   ├─> Logs success/failure
   │   └─> Decodes successful codes
   └─> Returns FuzzResults with valid IOCTLs

4. User calls: device.FormatIOCTL(0x220004)
   └─> Returns: "IOCTL 0x220004 [Device Type: 0x22, Function: 1, Method: BUFFERED, Access: ANY]"
```

### Workflow 2: Hook Process and Capture IOCTLs

```
1. User calls: inject.InjectDLL(1234, "C:\\hookdll.dll")
   ├─> Opens target process (PID 1234)
   ├─> Allocates memory in remote process
   ├─> Writes DLL path to remote memory
   ├─> Creates remote thread to call LoadLibraryA
   └─> DLL is loaded into target process

2. DLL executes: DllMain(DLL_PROCESS_ATTACH)
   ├─> Initializes IPC (named pipe to host)
   ├─> Installs IAT hook on kernel32!DeviceIoControl
   ├─> Replaces IAT entry with DeviceIoControlDetour
   └─> Creates trampoline for original function

3. Target app calls: DeviceIoControl(...)
   ├─> CPU jumps to DeviceIoControlDetour (our hook)
   ├─> Hook logs: device, IOCTL code, buffers
   ├─> Sends capture via named pipe to host
   ├─> Calls original DeviceIoControl via trampoline
   └─> Returns result to target app

4. Host process receives capture:
   ├─> capture.Session writes to C:\Captures\001.conf
   ├─> capture.Session writes to C:\Captures\001.data
   └─> User can replay later with ReplayIOCTL()
```

### Workflow 3: System-Wide Monitoring with ETW

```
1. User calls: etw.NewSession("MyIOCTLTrace")
   ├─> Creates ETW trace session
   ├─> Requires administrator privileges
   └─> Returns Session object

2. User calls: session.EnableProvider(KERNEL_FILE_GUID, level, keywords)
   ├─> Enables kernel file I/O provider
   ├─> Kernel starts logging I/O events
   └─> Events queued for consumption

3. User calls: consumer.StartConsuming(callback)
   ├─> Opens trace for real-time processing
   ├─> Processes events in background goroutine
   ├─> For each event:
   │   ├─> Parses EVENT_RECORD
   │   ├─> Extracts I/O details (if IOCTL-related)
   │   └─> Calls user callback
   └─> User callback logs/analyzes IOCTLs

4. User calls: session.Stop()
   ├─> Stops ETW trace
   └─> Cleanup resources
```

### Workflow 4: Replay Captured IOCTL

```
1. User has files: C:\Captures\001.conf, C:\Captures\001.data

2. User calls: capture.LoadCaptureConfig("C:\\Captures\\001.conf")
   ├─> Parses .conf file
   ├─> Returns CaptureConfig:
   │   ├─> DevicePath: \\.\MyDevice
   │   ├─> IOCTLCode: 0x220004
   │   ├─> InputSize: 64
   │   └─> OutputSize: 256
   └─> Decodes IOCTL code

3. User calls: capture.LoadCaptureData("C:\\Captures\\001.data")
   ├─> Reads binary .data file
   └─> Returns CaptureData with buffers

4. User calls: capture.ReplayIOCTL(configPath, dataPath)
   ├─> Opens device from config
   ├─> Calls DeviceIoControl with:
   │   ├─> IOCTL code from config
   │   ├─> Input buffer from data
   │   └─> Output buffer size from config
   ├─> Receives response
   └─> Returns IOCTLResponse with result
```

---

## Implementation Priorities

### Phase 1: Core Extensions (Week 1-2)
**Goal**: Essential functionality without hooking

Files to implement:
1. ✅ `device/decoder.go` - IOCTL decoding
2. ✅ `device/known_ioctls.go` - IOCTL database
3. ✅ `device/discovery.go` - Device delta detection
4. ✅ `device/capture.go` - .conf/.data file format
5. ✅ `service/driver_query.go` - Driver status queries

**Deliverable**: Can load drivers, enumerate devices, decode IOCTLs, save/replay captures manually

### Phase 2: Hooking Framework (Week 3-4)
**Goal**: User-mode hooking capabilities

Files to implement:
1. ✅ `memory/process_memory.go`
2. ✅ `memory/protection.go`
3. ✅ `memory/allocation.go`
4. ✅ `pe/parser.go`
5. ✅ `pe/imports.go`
6. ✅ `hook/iat_hook.go`
7. ✅ `asm/x64_encoder.go`
8. ✅ `asm/jump_gen.go`
9. ✅ `hook/inline_hook.go`
10. ✅ `hook/trampoline.go`
11. ✅ `hook/hook_manager.go`

**Deliverable**: Can hook DeviceIoControl in current process

### Phase 3: Injection Framework (Week 5-6)
**Goal**: Inject into remote processes

Files to implement:
1. ✅ `inject/dll_inject.go`
2. ✅ `inject/hookdll/main.go`
3. ✅ `inject/hookdll/hooks.go`
4. ✅ `ipc/named_pipe.go`
5. ✅ `capture/hook_bridge.go`
6. ✅ `capture/session.go`

**Deliverable**: Can inject into target process and capture IOCTLs to files

### Phase 4: Advanced Features (Week 7-8)
**Goal**: Fuzzing, ETW, WMI

Files to implement:
1. ✅ `device/fuzzer.go`
2. ✅ `wmi/query.go`
3. ✅ `wmi/driver_info.go`
4. ✅ `wmi/device_info.go`
5. ✅ `etw/session.go` (optional)
6. ✅ `etw/consumer.go` (optional)

**Deliverable**: Can fuzz drivers, query system via WMI, optional ETW monitoring

### Phase 5: Tools & Examples (Week 9-10)
**Goal**: User-facing tools

Files to implement:
1. ✅ `tools/winxctl/main.go` - CLI tool
2. ✅ `examples/load_driver/main.go`
3. ✅ `examples/hook_process/main.go`
4. ✅ `examples/capture_ioctls/main.go`
5. ✅ `examples/fuzz_driver/main.go`

**Deliverable**: Complete CLI tool and examples

---

### Phase 6: Windows Artifact Parsers (Week 11-14)
**Goal**: Core Windows artifact parsing capabilities

Files to implement:
1. ✅ `internal/artifacts/prefetch/parser.go` - Prefetch parser (pf)
2. ✅ `internal/artifacts/prefetch/decompress.go` - MAM decompression
3. ✅ `internal/artifacts/lnk/parser.go` - LNK parser (lp)
4. ✅ `internal/artifacts/jumplist/parser.go` - Jump list parser (jmp)
5. ✅ `internal/artifacts/jumplist/olecf.go` - OLE Compound File parsing
6. ✅ `device/usb.go` - USB storage parser (usp)
7. ✅ `internal/artifacts/shellbag/parser.go` - ShellBag parser (sbag)
8. ✅ `internal/artifacts/shimcache/parser.go` - ShimCache parser (wacu)
9. ✅ `internal/artifacts/shimdb/parser.go` - SDB parser (shims)
10. ✅ `internal/artifacts/activitiescache/parser.go` - ActivitiesCache (tac)
11. ✅ `internal/artifacts/indexdat/parser.go` - index.dat parser (id)
12. ✅ `internal/artifacts/recycle/parser.go` - Recycle Bin parser (tia)
13. ✅ `internal/artifacts/wpn/parser.go` - WPN database parser (wpn)
14. ✅ `internal/artifacts/backstage/parser.go` - Office Backstage parser (bs)

**Deliverable**: Parse all major Windows artifacts for timeline analysis

---

### Phase 7: Browser Artifacts (Week 15-16)
**Goal**: Browser history, cookies, cache parsing

Files to implement:
1. ✅ `internal/browser/chromium/parser.go` - Chromium parser (csp)
2. ✅ `internal/browser/chromium/history.go` - History database
3. ✅ `internal/browser/chromium/cookies.go` - Cookies database
4. ✅ `internal/browser/chromium/cache/parser.go` - Chromium cache (ccp)
5. ✅ `internal/browser/mozilla/parser.go` - Mozilla parser (msp)
6. ✅ `internal/browser/mozilla/places.go` - places.sqlite
7. ✅ `internal/browser/mozilla/cache/parser.go` - Mozilla cache (mcp)
8. ✅ `internal/browser/safari/parser.go` - Safari parser (sap)

**Deliverable**: Parse all major browser artifacts

---

### Phase 8: Registry & Event Log Analysis (Week 17-20)
**Goal**: Offline registry and event log parsing

Files to implement:
1. ✅ `registry/hive.go` - Registry hive parser (yaru)
2. ✅ `registry/cell.go` - Cell parsing
3. ✅ `registry/value.go` - Value data parsing
4. ✅ `registry/dirty.go` - Transaction log parsing
5. ✅ `internal/evtx/parser.go` - EVTX parser (evtwalk/evtx_view)
6. ✅ `internal/evtx/chunk.go` - Chunk parsing
7. ✅ `internal/evtx/binxml.go` - Binary XML parsing
8. ✅ `internal/evtx/fragment/parser.go` - Fragment recovery (evtfx)
9. ✅ `internal/messagetable/parser.go` - Message table parser (elmo)
10. ✅ `internal/cafae/extractor.go` - Account artifact extractor (cafae)
11. ✅ `internal/tela/analyzer.go` - Event correlation (tela)

**Deliverable**: Complete registry and event log analysis capabilities

---

### Phase 9: NTFS Filesystem Analysis (Week 21-26)
**Goal**: Complete NTFS parsing and analysis

Files to implement:
1. ✅ `internal/filesystem/ntfs/volume.go` - NTFS volume handling
2. ✅ `internal/filesystem/ntfs/boot.go` - Boot sector parsing
3. ✅ `internal/filesystem/mft/parser.go` - $MFT parser (ntfswalk)
4. ✅ `internal/filesystem/mft/attribute.go` - Attribute parsing
5. ✅ `internal/filesystem/mft/runlist.go` - Data run parsing
6. ✅ `internal/filesystem/usnjrnl/parser.go` - $UsnJrnl parser (jp)
7. ✅ `internal/filesystem/logfile/parser.go` - $LogFile parser (mala)
8. ✅ `internal/filesystem/indx/parser.go` - INDX slack parser (wisp)
9. ✅ `internal/filesystem/indx/carver.go` - Deleted entry recovery
10. ✅ `internal/filesystem/ntfsdir/enumerator.go` - Directory enum (ntfsdir)
11. ✅ `internal/filesystem/ntfscopy/copy.go` - NTFS file copy (ntfscopy)
12. ✅ `internal/filesystem/gena/engine.go` - NTFS analysis engine (gena)
13. ✅ `internal/filesystem/fat/fat32.go` - FAT32 parser (fata)
14. ✅ `internal/filesystem/fat/exfat.go` - exFAT parser

**Deliverable**: Complete filesystem analysis for NTFS and FAT

---

### Phase 10: Network & PE Utilities (Week 27-30)
**Goal**: Network tools and extended PE analysis

Files to implement:
1. ✅ `internal/network/dns/query.go` - DNS query utility (dqu)
2. ✅ `internal/network/pcap/parser.go` - PCAP parser (pic)
3. ✅ `internal/network/pcap/icmp.go` - ICMP carving
4. ✅ `internal/network/netxfer/server.go` - Network transfer server (nx)
5. ✅ `internal/network/netxfer/client.go` - Network transfer client
6. ✅ `internal/network/minx/agent.go` - MINX agent (minx)
7. ✅ `internal/pe/scanner.go` - PE anomaly scanner (pescan)
8. ✅ `internal/pe/viewer.go` - PE viewer (pe_view)
9. ✅ `internal/pe/resources.go` - Resource parsing
10. ✅ `internal/pe/debug.go` - Debug directory

**Deliverable**: Complete network utilities and PE analysis

---

### Phase 11: Miscellaneous Utilities (Week 31-34)
**Goal**: VSS, symbols, disk utilities

Files to implement:
1. ✅ `internal/vss/enumerator.go` - VSS enumeration (vssenum)
2. ✅ `internal/vss/snapshot.go` - Snapshot access
3. ✅ `internal/vss/copy.go` - Copy from snapshots
4. ✅ `internal/symbol/fetch.go` - Symbol fetcher (sf)
5. ✅ `internal/symbol/pdb.go` - PDB parsing
6. ✅ `internal/csvdx/parser.go` - CSV utilities (csvdx)
7. ✅ `internal/disk/disk.go` - Raw disk access (dup)
8. ✅ `internal/disk/partition.go` - Partition parsing
9. ✅ `internal/disk/imaging.go` - Forensic imaging

**Deliverable**: Complete miscellaneous utilities

---

### Phase 12: Example Programs (Week 35-36)
**Goal**: Comprehensive example programs demonstrating all capabilities

Files to implement:
1. ✅ `examples/parse_prefetch/main.go`
2. ✅ `examples/parse_lnk/main.go`
3. ✅ `examples/parse_jumplist/main.go`
4. ✅ `examples/parse_usb/main.go`
5. ✅ `examples/parse_evtx/main.go`
6. ✅ `examples/parse_mft/main.go`
7. ✅ `examples/parse_registry/main.go`
8. ✅ `examples/parse_chromium/main.go`
9. ✅ `examples/pe_scan/main.go`
10. ✅ `examples/vss_enumerate/main.go`

**Deliverable**: Complete example programs for all major features

---

## Summary

**Total New Functions**: ~500+ functions across 40+ packages
**Total New Files**: ~150+ files
**Dependencies**: Only `golang.org/x/sys/windows` (for Windows API syscalls)
**Approach**: Build everything from scratch in pure Go
**Estimated Timeline**: 36 weeks for complete implementation

This design provides a comprehensive Windows analysis framework including:

### IOCTL++ Capabilities
- User-mode API hooking (IAT + inline)
- Process injection (DLL, reflective, thread hijack)
- ETW monitoring
- WMI queries
- Custom capture/replay format
- IOCTL fuzzing engine
- Complete assembly framework for hooking

### Artifact Analysis
- **Prefetch Parser (pf)** - Execution timeline from prefetch files
- **LNK Parser (lp)** - Shell link file analysis
- **Jump List Parser (jmp)** - Recent/frequent document tracking
- **USB Storage Parser (usp)** - USB device connection history
- **ShellBag Parser (sbag)** - Folder access history
- **AppCompat Cache (wacu)** - Application execution evidence
- **Shim Database (shims)** - Application compatibility shims
- **ActivitiesCache (tac)** - Windows Timeline activities
- **index.dat Parser (id)** - Legacy IE history
- **Recycle Bin (tia)** - Deleted file analysis
- **WPN Database (wpn)** - Push notification history
- **Office Backstage (bs)** - Recent Office documents

### Browser Artifacts
- **Chromium Parser (csp)** - Chrome/Edge/Brave history, cookies, downloads
- **Chromium Cache (ccp)** - Browser cache analysis
- **Mozilla Parser (msp)** - Firefox history, cookies, forms
- **Mozilla Cache (mcp)** - Firefox cache analysis
- **Safari Parser (sap)** - Safari history and bookmarks

### Registry & Event Log Analysis
- **Registry Utility (yaru)** - Offline registry hive parsing
- **EVTX Parser (evtwalk/evtx_view)** - Windows event log parsing
- **EVTX Fragments (evtfx)** - Fragment recovery and carving
- **Message Tables (elmo)** - Offline event message resolution
- **Account Artifacts (cafae)** - SAM/SECURITY/SOFTWARE extraction
- **Event Analysis (tela)** - Event correlation and timeline

### NTFS Filesystem Analysis
- **$MFT Parser (ntfswalk)** - Master File Table analysis
- **$UsnJrnl Parser (jp)** - Change journal parsing
- **$LogFile Parser (mala)** - Transaction log analysis
- **INDX Slack (wisp)** - Deleted file recovery from indexes
- **NTFS Directory (ntfsdir)** - Raw directory enumeration
- **NTFS Copy (ntfscopy)** - Copy locked/in-use files
- **NTFS Analysis (gena)** - Graphical analysis engine
- **FAT Analysis (fata)** - FAT32/exFAT parsing

### Network Utilities
- **DNS Query (dqu)** - DNS queries and cache inspection
- **PCAP/ICMP Carver (pic)** - Packet capture analysis
- **Network Transfer (nx)** - Secure file transfer
- **MINX Agent (minx)** - Modular inspection agent

### PE & Miscellaneous
- **PE Viewer (pe_view)** - Complete PE file analysis
- **PE Scanner (pescan)** - Anomaly and packer detection
- **VSS Enumerator (vssenum)** - Shadow copy analysis
- **Symbol Fetch (sf)** - Microsoft symbol server client
- **CSV Utilities (csvdx)** - Data exchange and transformation
- **Disk Utility (dup)** - Raw disk access and forensic imaging

All functionality is modular and can be used independently or combined for comprehensive Windows system analysis and research.
