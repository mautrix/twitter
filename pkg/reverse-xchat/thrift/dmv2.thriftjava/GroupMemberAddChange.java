package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.clientauthutils.C0026b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000@\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0010\u000e\n\u0002\b\u0006\n\u0002\u0010\t\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0010\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\b\b\u0086\b\u0018\u0000 -2\u00020\u0001:\u0002.-Bo\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\u000e\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\u000e\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0003\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0003\u0012\b\u0010\t\u001a\u0004\u0018\u00010\u0003\u0012\b\u0010\u000b\u001a\u0004\u0018\u00010\n\u0012\u000e\u0010\f\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002¢\u0006\u0004\b\r\u0010\u000eJ\u0017\u0010\u0012\u001a\u00020\u00112\u0006\u0010\u0010\u001a\u00020\u000fH\u0016¢\u0006\u0004\b\u0012\u0010\u0013J\u0018\u0010\u0014\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0014\u0010\u0015J\u0018\u0010\u0016\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0016\u0010\u0015J\u0018\u0010\u0017\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0017\u0010\u0015J\u0012\u0010\u0018\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u0018\u0010\u0019J\u0012\u0010\u001a\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u001a\u0010\u0019J\u0012\u0010\u001b\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u001b\u0010\u0019J\u0012\u0010\u001c\u001a\u0004\u0018\u00010\nHÆ\u0003¢\u0006\u0004\b\u001c\u0010\u001dJ\u0018\u0010\u001e\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001e\u0010\u0015J\u0088\u0001\u0010\u001f\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\u0010\b\u0002\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\u0010\b\u0002\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00032\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u00032\n\b\u0002\u0010\t\u001a\u0004\u0018\u00010\u00032\n\b\u0002\u0010\u000b\u001a\u0004\u0018\u00010\n2\u0010\b\u0002\u0010\f\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u001f\u0010 J\u0010\u0010!\u001a\u00020\u0003HÖ\u0001¢\u0006\u0004\b!\u0010\u0019J\u0010\u0010#\u001a\u00020\"HÖ\u0001¢\u0006\u0004\b#\u0010$J\u001a\u0010(\u001a\u00020'2\b\u0010&\u001a\u0004\u0018\u00010%HÖ\u0003¢\u0006\u0004\b(\u0010)R\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010*R\u001c\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010*R\u001c\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010*R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010+R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010+R\u0016\u0010\t\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\t\u0010+R\u0016\u0010\u000b\u001a\u0004\u0018\u00010\n8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000b\u0010,R\u001c\u0010\f\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\f\u0010*¨\u0006/"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "Lcom/bendb/thrifty/a;", "", "", "member_ids", "current_member_ids", "current_admin_ids", "current_title", "current_avatar_url", "conversation_key_version", "", "current_ttl_msec", "current_pending_member_ids", "<init>", "(Ljava/util/List;Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/Long;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "component3", "component4", "()Ljava/lang/String;", "component5", "component6", "component7", "()Ljava/lang/Long;", "component8", "copy", "(Ljava/util/List;Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/Long;Ljava/util/List;)Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Ljava/lang/String;", "Ljava/lang/Long;", "Companion", "GroupMemberAddChangeAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class GroupMemberAddChange implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String conversation_key_version;

    @JvmField
    @InterfaceC88465b
    public final List current_admin_ids;

    @JvmField
    @InterfaceC88465b
    public final String current_avatar_url;

    @JvmField
    @InterfaceC88465b
    public final List current_member_ids;

    @JvmField
    @InterfaceC88465b
    public final List current_pending_member_ids;

    @JvmField
    @InterfaceC88465b
    public final String current_title;

    @JvmField
    @InterfaceC88465b
    public final Long current_ttl_msec;

    @JvmField
    @InterfaceC88465b
    public final List member_ids;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new GroupMemberAddChangeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupMemberAddChange$GroupMemberAddChangeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/GroupMemberAddChange;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/GroupMemberAddChange;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class GroupMemberAddChangeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public GroupMemberAddChange m83723read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            ArrayList arrayList2 = null;
            ArrayList arrayList3 = null;
            String string = null;
            String string2 = null;
            String string3 = null;
            Long lValueOf = null;
            ArrayList arrayList4 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    int i = 0;
                    switch (c11265cMo14127V2.f38393b) {
                        case 1:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i2 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList5 = new ArrayList(i2);
                                while (i < i2) {
                                    arrayList5.add(protocol.readString());
                                    i++;
                                }
                                arrayList = arrayList5;
                                break;
                            }
                        case 2:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i3 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList6 = new ArrayList(i3);
                                while (i < i3) {
                                    arrayList6.add(protocol.readString());
                                    i++;
                                }
                                arrayList2 = arrayList6;
                                break;
                            }
                        case 3:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i4 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList7 = new ArrayList(i4);
                                while (i < i4) {
                                    arrayList7.add(protocol.readString());
                                    i++;
                                }
                                arrayList3 = arrayList7;
                                break;
                            }
                        case 4:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string = protocol.readString();
                                break;
                            }
                        case 5:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string2 = protocol.readString();
                                break;
                            }
                        case 6:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string3 = protocol.readString();
                                break;
                            }
                        case 7:
                            if (b != 10) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                lValueOf = Long.valueOf(protocol.mo14124H0());
                                break;
                            }
                        case 8:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i5 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList8 = new ArrayList(i5);
                                while (i < i5) {
                                    arrayList8.add(protocol.readString());
                                    i++;
                                }
                                arrayList4 = arrayList8;
                                break;
                            }
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                    }
                } else {
                    return new GroupMemberAddChange(arrayList, arrayList2, arrayList3, string, string2, string3, lValueOf, arrayList4);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a GroupMemberAddChange struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("GroupMemberAddChange");
            if (struct.member_ids != null) {
                protocol.mo14136v3("member_ids", 1, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.member_ids.size());
                Iterator it = struct.member_ids.iterator();
                while (it.hasNext()) {
                    protocol.mo14137w0((String) it.next());
                }
            }
            if (struct.current_member_ids != null) {
                protocol.mo14136v3("current_member_ids", 2, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.current_member_ids.size());
                Iterator it2 = struct.current_member_ids.iterator();
                while (it2.hasNext()) {
                    protocol.mo14137w0((String) it2.next());
                }
            }
            if (struct.current_admin_ids != null) {
                protocol.mo14136v3("current_admin_ids", 3, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.current_admin_ids.size());
                Iterator it3 = struct.current_admin_ids.iterator();
                while (it3.hasNext()) {
                    protocol.mo14137w0((String) it3.next());
                }
            }
            if (struct.current_title != null) {
                protocol.mo14136v3("current_title", 4, (byte) 11);
                protocol.mo14137w0(struct.current_title);
            }
            if (struct.current_avatar_url != null) {
                protocol.mo14136v3("current_avatar_url", 5, (byte) 11);
                protocol.mo14137w0(struct.current_avatar_url);
            }
            if (struct.conversation_key_version != null) {
                protocol.mo14136v3("conversation_key_version", 6, (byte) 11);
                protocol.mo14137w0(struct.conversation_key_version);
            }
            if (struct.current_ttl_msec != null) {
                protocol.mo14136v3("current_ttl_msec", 7, (byte) 10);
                protocol.mo14121B3(struct.current_ttl_msec.longValue());
            }
            if (struct.current_pending_member_ids != null) {
                protocol.mo14136v3("current_pending_member_ids", 8, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.current_pending_member_ids.size());
                Iterator it4 = struct.current_pending_member_ids.iterator();
                while (it4.hasNext()) {
                    protocol.mo14137w0((String) it4.next());
                }
            }
            protocol.mo14134i0();
        }
    }

    public GroupMemberAddChange(@InterfaceC88465b List list, @InterfaceC88465b List list2, @InterfaceC88465b List list3, @InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b Long l, @InterfaceC88465b List list4) {
        this.member_ids = list;
        this.current_member_ids = list2;
        this.current_admin_ids = list3;
        this.current_title = str;
        this.current_avatar_url = str2;
        this.conversation_key_version = str3;
        this.current_ttl_msec = l;
        this.current_pending_member_ids = list4;
    }

    public static /* synthetic */ GroupMemberAddChange copy$default(GroupMemberAddChange groupMemberAddChange, List list, List list2, List list3, String str, String str2, String str3, Long l, List list4, int i, Object obj) {
        return groupMemberAddChange.copy((i & 1) != 0 ? groupMemberAddChange.member_ids : list, (i & 2) != 0 ? groupMemberAddChange.current_member_ids : list2, (i & 4) != 0 ? groupMemberAddChange.current_admin_ids : list3, (i & 8) != 0 ? groupMemberAddChange.current_title : str, (i & 16) != 0 ? groupMemberAddChange.current_avatar_url : str2, (i & 32) != 0 ? groupMemberAddChange.conversation_key_version : str3, (i & 64) != 0 ? groupMemberAddChange.current_ttl_msec : l, (i & 128) != 0 ? groupMemberAddChange.current_pending_member_ids : list4);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getMember_ids() {
        return this.member_ids;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getCurrent_member_ids() {
        return this.current_member_ids;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final List getCurrent_admin_ids() {
        return this.current_admin_ids;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getCurrent_title() {
        return this.current_title;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getCurrent_avatar_url() {
        return this.current_avatar_url;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final String getConversation_key_version() {
        return this.conversation_key_version;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final Long getCurrent_ttl_msec() {
        return this.current_ttl_msec;
    }

    @InterfaceC88465b
    /* renamed from: component8, reason: from getter */
    public final List getCurrent_pending_member_ids() {
        return this.current_pending_member_ids;
    }

    @InterfaceC88464a
    public final GroupMemberAddChange copy(@InterfaceC88465b List member_ids, @InterfaceC88465b List current_member_ids, @InterfaceC88465b List current_admin_ids, @InterfaceC88465b String current_title, @InterfaceC88465b String current_avatar_url, @InterfaceC88465b String conversation_key_version, @InterfaceC88465b Long current_ttl_msec, @InterfaceC88465b List current_pending_member_ids) {
        return new GroupMemberAddChange(member_ids, current_member_ids, current_admin_ids, current_title, current_avatar_url, conversation_key_version, current_ttl_msec, current_pending_member_ids);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof GroupMemberAddChange)) {
            return false;
        }
        GroupMemberAddChange groupMemberAddChange = (GroupMemberAddChange) other;
        return Intrinsics.m65267c(this.member_ids, groupMemberAddChange.member_ids) && Intrinsics.m65267c(this.current_member_ids, groupMemberAddChange.current_member_ids) && Intrinsics.m65267c(this.current_admin_ids, groupMemberAddChange.current_admin_ids) && Intrinsics.m65267c(this.current_title, groupMemberAddChange.current_title) && Intrinsics.m65267c(this.current_avatar_url, groupMemberAddChange.current_avatar_url) && Intrinsics.m65267c(this.conversation_key_version, groupMemberAddChange.conversation_key_version) && Intrinsics.m65267c(this.current_ttl_msec, groupMemberAddChange.current_ttl_msec) && Intrinsics.m65267c(this.current_pending_member_ids, groupMemberAddChange.current_pending_member_ids);
    }

    public int hashCode() {
        List list = this.member_ids;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        List list2 = this.current_member_ids;
        int iHashCode2 = (iHashCode + (list2 == null ? 0 : list2.hashCode())) * 31;
        List list3 = this.current_admin_ids;
        int iHashCode3 = (iHashCode2 + (list3 == null ? 0 : list3.hashCode())) * 31;
        String str = this.current_title;
        int iHashCode4 = (iHashCode3 + (str == null ? 0 : str.hashCode())) * 31;
        String str2 = this.current_avatar_url;
        int iHashCode5 = (iHashCode4 + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.conversation_key_version;
        int iHashCode6 = (iHashCode5 + (str3 == null ? 0 : str3.hashCode())) * 31;
        Long l = this.current_ttl_msec;
        int iHashCode7 = (iHashCode6 + (l == null ? 0 : l.hashCode())) * 31;
        List list4 = this.current_pending_member_ids;
        return iHashCode7 + (list4 != null ? list4.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        List list = this.member_ids;
        List list2 = this.current_member_ids;
        List list3 = this.current_admin_ids;
        String str = this.current_title;
        String str2 = this.current_avatar_url;
        String str3 = this.conversation_key_version;
        Long l = this.current_ttl_msec;
        List list4 = this.current_pending_member_ids;
        StringBuilder sb = new StringBuilder("GroupMemberAddChange(member_ids=");
        sb.append(list);
        sb.append(", current_member_ids=");
        sb.append(list2);
        sb.append(", current_admin_ids=");
        sb.append(list3);
        sb.append(", current_title=");
        sb.append(str);
        sb.append(", current_avatar_url=");
        C0026b.m37b(sb, str2, ", conversation_key_version=", str3, ", current_ttl_msec=");
        sb.append(l);
        sb.append(", current_pending_member_ids=");
        sb.append(list4);
        sb.append(Separators.RPAREN);
        return sb.toString();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}