package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000L\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0010\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u000e\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\u0013\u0014\u0015\u0016\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\f\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f !\"¨\u0006#"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "GroupCreate", "GroupTitleChange", "GroupAvatarChange", "GroupAdminAdd", "GroupMemberAdd", "GroupAdminRemove", "GroupMemberRemove", "GroupInviteEnable", "GroupInviteDisable", "GroupJoinRequest", "GroupJoinReject", "Unknown", "GroupChangeAdapter", "Lcom/x/dmv2/thriftjava/GroupChange$GroupAdminAdd;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupAdminRemove;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupAvatarChange;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupCreate;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupInviteDisable;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupInviteEnable;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupJoinReject;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupJoinRequest;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupMemberAdd;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupMemberRemove;", "Lcom/x/dmv2/thriftjava/GroupChange$GroupTitleChange;", "Lcom/x/dmv2/thriftjava/GroupChange$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class GroupChange implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new GroupChangeAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupAdminAdd;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupAdminAddChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupAdminAddChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupAdminAddChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupAdminAdd extends GroupChange {

        @InterfaceC88464a
        private final GroupAdminAddChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupAdminAdd(@InterfaceC88464a GroupAdminAddChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupAdminAdd copy$default(GroupAdminAdd groupAdminAdd, GroupAdminAddChange groupAdminAddChange, int i, Object obj) {
            if ((i & 1) != 0) {
                groupAdminAddChange = groupAdminAdd.value;
            }
            return groupAdminAdd.copy(groupAdminAddChange);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final GroupAdminAddChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupAdminAdd copy(@InterfaceC88464a GroupAdminAddChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupAdminAdd(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupAdminAdd) && Intrinsics.m65267c(this.value, ((GroupAdminAdd) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final GroupAdminAddChange m76747getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_admin_add=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupAdminRemove;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupAdminRemoveChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupAdminRemoveChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupAdminRemoveChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupAdminRemove extends GroupChange {

        @InterfaceC88464a
        private final GroupAdminRemoveChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupAdminRemove(@InterfaceC88464a GroupAdminRemoveChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupAdminRemove copy$default(GroupAdminRemove groupAdminRemove, GroupAdminRemoveChange groupAdminRemoveChange, int i, Object obj) {
            if ((i & 1) != 0) {
                groupAdminRemoveChange = groupAdminRemove.value;
            }
            return groupAdminRemove.copy(groupAdminRemoveChange);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final GroupAdminRemoveChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupAdminRemove copy(@InterfaceC88464a GroupAdminRemoveChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupAdminRemove(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupAdminRemove) && Intrinsics.m65267c(this.value, ((GroupAdminRemove) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final GroupAdminRemoveChange m76748getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_admin_remove=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupAvatarChange;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupAvatarUrlChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupAvatarUrlChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupAvatarUrlChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupAvatarChange extends GroupChange {

        @InterfaceC88464a
        private final GroupAvatarUrlChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupAvatarChange(@InterfaceC88464a GroupAvatarUrlChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupAvatarChange copy$default(GroupAvatarChange groupAvatarChange, GroupAvatarUrlChange groupAvatarUrlChange, int i, Object obj) {
            if ((i & 1) != 0) {
                groupAvatarUrlChange = groupAvatarChange.value;
            }
            return groupAvatarChange.copy(groupAvatarUrlChange);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final GroupAvatarUrlChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupAvatarChange copy(@InterfaceC88464a GroupAvatarUrlChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupAvatarChange(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupAvatarChange) && Intrinsics.m65267c(this.value, ((GroupAvatarChange) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final GroupAvatarUrlChange m76749getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_avatar_change=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupChangeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/GroupChange;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/GroupChange;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/GroupChange;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class GroupChangeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public GroupChange m85645read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            GroupChange groupCreate;
            Intrinsics.m65272h(protocol, "protocol");
            GroupChange groupChange = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (groupChange != null) {
                        return groupChange;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            groupCreate = new GroupCreate((com.x.dmv2.thriftjava.GroupCreate) com.x.dmv2.thriftjava.GroupCreate.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                        if (b == 12) {
                            groupCreate = new GroupTitleChange((com.x.dmv2.thriftjava.GroupTitleChange) com.x.dmv2.thriftjava.GroupTitleChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 3:
                        if (b == 12) {
                            groupCreate = new GroupAvatarChange((GroupAvatarUrlChange) GroupAvatarUrlChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            groupCreate = new GroupAdminAdd((GroupAdminAddChange) GroupAdminAddChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            groupCreate = new GroupMemberAdd((GroupMemberAddChange) GroupMemberAddChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            groupCreate = new GroupAdminRemove((GroupAdminRemoveChange) GroupAdminRemoveChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 7:
                        if (b == 12) {
                            groupCreate = new GroupMemberRemove((GroupMemberRemoveChange) GroupMemberRemoveChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 8:
                        if (b == 12) {
                            groupCreate = new GroupInviteEnable((com.x.dmv2.thriftjava.GroupInviteEnable) com.x.dmv2.thriftjava.GroupInviteEnable.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 9:
                        if (b == 12) {
                            groupCreate = new GroupInviteDisable((com.x.dmv2.thriftjava.GroupInviteDisable) com.x.dmv2.thriftjava.GroupInviteDisable.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 10:
                        if (b == 12) {
                            groupCreate = new GroupJoinRequest((com.x.dmv2.thriftjava.GroupJoinRequest) com.x.dmv2.thriftjava.GroupJoinRequest.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 11:
                        if (b == 12) {
                            groupCreate = new GroupJoinReject((com.x.dmv2.thriftjava.GroupJoinReject) com.x.dmv2.thriftjava.GroupJoinReject.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    default:
                        groupChange = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                }
                groupChange = groupCreate;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a GroupChange struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("GroupChange");
            if (struct instanceof GroupCreate) {
                protocol.mo14136v3("group_create", 1, (byte) 12);
                com.x.dmv2.thriftjava.GroupCreate.ADAPTER.write(protocol, ((GroupCreate) struct).m76750getValue());
            } else if (struct instanceof GroupTitleChange) {
                protocol.mo14136v3("group_title_change", 2, (byte) 12);
                com.x.dmv2.thriftjava.GroupTitleChange.ADAPTER.write(protocol, ((GroupTitleChange) struct).m76757getValue());
            } else if (struct instanceof GroupAvatarChange) {
                protocol.mo14136v3("group_avatar_change", 3, (byte) 12);
                GroupAvatarUrlChange.ADAPTER.write(protocol, ((GroupAvatarChange) struct).m76749getValue());
            } else if (struct instanceof GroupAdminAdd) {
                protocol.mo14136v3("group_admin_add", 4, (byte) 12);
                GroupAdminAddChange.ADAPTER.write(protocol, ((GroupAdminAdd) struct).m76747getValue());
            } else if (struct instanceof GroupMemberAdd) {
                protocol.mo14136v3("group_member_add", 5, (byte) 12);
                GroupMemberAddChange.ADAPTER.write(protocol, ((GroupMemberAdd) struct).m76755getValue());
            } else if (struct instanceof GroupAdminRemove) {
                protocol.mo14136v3("group_admin_remove", 6, (byte) 12);
                GroupAdminRemoveChange.ADAPTER.write(protocol, ((GroupAdminRemove) struct).m76748getValue());
            } else if (struct instanceof GroupMemberRemove) {
                protocol.mo14136v3("group_member_remove", 7, (byte) 12);
                GroupMemberRemoveChange.ADAPTER.write(protocol, ((GroupMemberRemove) struct).m76756getValue());
            } else if (struct instanceof GroupInviteEnable) {
                protocol.mo14136v3("group_invite_enable", 8, (byte) 12);
                com.x.dmv2.thriftjava.GroupInviteEnable.ADAPTER.write(protocol, ((GroupInviteEnable) struct).m76752getValue());
            } else if (struct instanceof GroupInviteDisable) {
                protocol.mo14136v3("group_invite_disable", 9, (byte) 12);
                com.x.dmv2.thriftjava.GroupInviteDisable.ADAPTER.write(protocol, ((GroupInviteDisable) struct).m76751getValue());
            } else if (struct instanceof GroupJoinRequest) {
                protocol.mo14136v3("group_join_request", 10, (byte) 12);
                com.x.dmv2.thriftjava.GroupJoinRequest.ADAPTER.write(protocol, ((GroupJoinRequest) struct).m76754getValue());
            } else if (struct instanceof GroupJoinReject) {
                protocol.mo14136v3("group_join_reject", 11, (byte) 12);
                com.x.dmv2.thriftjava.GroupJoinReject.ADAPTER.write(protocol, ((GroupJoinReject) struct).m76753getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupCreate;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupCreate;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupCreate;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupCreate;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupCreate extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupCreate value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupCreate(@InterfaceC88464a com.x.dmv2.thriftjava.GroupCreate value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupCreate copy$default(GroupCreate groupCreate, com.x.dmv2.thriftjava.GroupCreate groupCreate2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupCreate2 = groupCreate.value;
            }
            return groupCreate.copy(groupCreate2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupCreate getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupCreate copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupCreate value) {
            Intrinsics.m65272h(value, "value");
            return new GroupCreate(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupCreate) && Intrinsics.m65267c(this.value, ((GroupCreate) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupCreate m76750getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_create=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupInviteDisable;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupInviteDisable;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupInviteDisable;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupInviteDisable;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupInviteDisable extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupInviteDisable value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupInviteDisable(@InterfaceC88464a com.x.dmv2.thriftjava.GroupInviteDisable value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupInviteDisable copy$default(GroupInviteDisable groupInviteDisable, com.x.dmv2.thriftjava.GroupInviteDisable groupInviteDisable2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupInviteDisable2 = groupInviteDisable.value;
            }
            return groupInviteDisable.copy(groupInviteDisable2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupInviteDisable getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupInviteDisable copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupInviteDisable value) {
            Intrinsics.m65272h(value, "value");
            return new GroupInviteDisable(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupInviteDisable) && Intrinsics.m65267c(this.value, ((GroupInviteDisable) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupInviteDisable m76751getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_invite_disable=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupInviteEnable;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupInviteEnable;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupInviteEnable extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupInviteEnable value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupInviteEnable(@InterfaceC88464a com.x.dmv2.thriftjava.GroupInviteEnable value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupInviteEnable copy$default(GroupInviteEnable groupInviteEnable, com.x.dmv2.thriftjava.GroupInviteEnable groupInviteEnable2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupInviteEnable2 = groupInviteEnable.value;
            }
            return groupInviteEnable.copy(groupInviteEnable2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupInviteEnable getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupInviteEnable copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupInviteEnable value) {
            Intrinsics.m65272h(value, "value");
            return new GroupInviteEnable(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupInviteEnable) && Intrinsics.m65267c(this.value, ((GroupInviteEnable) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupInviteEnable m76752getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_invite_enable=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupJoinReject;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupJoinReject;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupJoinReject;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupJoinReject;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupJoinReject extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupJoinReject value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupJoinReject(@InterfaceC88464a com.x.dmv2.thriftjava.GroupJoinReject value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupJoinReject copy$default(GroupJoinReject groupJoinReject, com.x.dmv2.thriftjava.GroupJoinReject groupJoinReject2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupJoinReject2 = groupJoinReject.value;
            }
            return groupJoinReject.copy(groupJoinReject2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupJoinReject getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupJoinReject copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupJoinReject value) {
            Intrinsics.m65272h(value, "value");
            return new GroupJoinReject(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupJoinReject) && Intrinsics.m65267c(this.value, ((GroupJoinReject) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupJoinReject m76753getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_join_reject=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupJoinRequest;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupJoinRequest;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupJoinRequest;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupJoinRequest;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupJoinRequest extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupJoinRequest value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupJoinRequest(@InterfaceC88464a com.x.dmv2.thriftjava.GroupJoinRequest value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupJoinRequest copy$default(GroupJoinRequest groupJoinRequest, com.x.dmv2.thriftjava.GroupJoinRequest groupJoinRequest2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupJoinRequest2 = groupJoinRequest.value;
            }
            return groupJoinRequest.copy(groupJoinRequest2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupJoinRequest getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupJoinRequest copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupJoinRequest value) {
            Intrinsics.m65272h(value, "value");
            return new GroupJoinRequest(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupJoinRequest) && Intrinsics.m65267c(this.value, ((GroupJoinRequest) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupJoinRequest m76754getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_join_request=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupMemberAdd;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupMemberAddChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupMemberAdd extends GroupChange {

        @InterfaceC88464a
        private final GroupMemberAddChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupMemberAdd(@InterfaceC88464a GroupMemberAddChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupMemberAdd copy$default(GroupMemberAdd groupMemberAdd, GroupMemberAddChange groupMemberAddChange, int i, Object obj) {
            if ((i & 1) != 0) {
                groupMemberAddChange = groupMemberAdd.value;
            }
            return groupMemberAdd.copy(groupMemberAddChange);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final GroupMemberAddChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupMemberAdd copy(@InterfaceC88464a GroupMemberAddChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupMemberAdd(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupMemberAdd) && Intrinsics.m65267c(this.value, ((GroupMemberAdd) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final GroupMemberAddChange m76755getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_member_add=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupMemberRemove;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupMemberRemoveChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupMemberRemoveChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupMemberRemoveChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupMemberRemove extends GroupChange {

        @InterfaceC88464a
        private final GroupMemberRemoveChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupMemberRemove(@InterfaceC88464a GroupMemberRemoveChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupMemberRemove copy$default(GroupMemberRemove groupMemberRemove, GroupMemberRemoveChange groupMemberRemoveChange, int i, Object obj) {
            if ((i & 1) != 0) {
                groupMemberRemoveChange = groupMemberRemove.value;
            }
            return groupMemberRemove.copy(groupMemberRemoveChange);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final GroupMemberRemoveChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupMemberRemove copy(@InterfaceC88464a GroupMemberRemoveChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupMemberRemove(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupMemberRemove) && Intrinsics.m65267c(this.value, ((GroupMemberRemove) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final GroupMemberRemoveChange m76756getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_member_remove=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$GroupTitleChange;", "Lcom/x/dmv2/thriftjava/GroupChange;", "value", "Lcom/x/dmv2/thriftjava/GroupTitleChange;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupTitleChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupTitleChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupTitleChange extends GroupChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupTitleChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupTitleChange(@InterfaceC88464a com.x.dmv2.thriftjava.GroupTitleChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupTitleChange copy$default(GroupTitleChange groupTitleChange, com.x.dmv2.thriftjava.GroupTitleChange groupTitleChange2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupTitleChange2 = groupTitleChange.value;
            }
            return groupTitleChange.copy(groupTitleChange2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupTitleChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupTitleChange copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupTitleChange value) {
            Intrinsics.m65272h(value, "value");
            return new GroupTitleChange(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupTitleChange) && Intrinsics.m65267c(this.value, ((GroupTitleChange) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupTitleChange m76757getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "GroupChange(group_title_change=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupChange$Unknown;", "Lcom/x/dmv2/thriftjava/GroupChange;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends GroupChange {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -648393154;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ GroupChange(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private GroupChange() {
    }
}
